provider "azurerm" {
  subscription_id = "${var.azure_subscription_id}"
  client_id       = "${var.azure_client_id}"
  client_secret   = "${var.azure_client_secret}"
  tenant_id       = "${var.azure_tenant_id}"
}

resource "random_string" "resource_group_name" {
  length  = 8
  special = false
}

resource "random_string" "net_name" {
  length  = 8
  special = false
}

# resource "ansible_host" "hideNsneak" {
#   count = "${var.azure_instance_count}"

#   //Element
#   inventory_hostname = "${azurerm_public_ip.public_ip.*.ip_address[count.index]}"
#   groups             = "${var.ansible_groups}"

#   vars {
#     ansible_user       = "${var.azure_default_username}"
#     ansible_connection = "ssh"

#     ansible_ssh_private_key_file = "${var.azure_private_key_file}"
#     ansible_ssh_common_args      = "-o StrictHostKeyChecking=no"
#   }

#   depends_on = ["azure_virtual_machine.hideNsneak"]
# }

//TODO: Resource group may not need to be created with each module
resource "azurerm_resource_group" "hideNsneak" {
  name     = "hideNsneak${random_string.resource_group_name.result}"
  count    = 1
  location = "${var.azure_location}"
}

//TODO: Figure how to map multiple IPs to instances
resource "azurerm_public_ip" "public_ip" {
  count                        = "${var.azure_instance_count}"
  name                         = "hideNsneak"
  location                     = "${azurerm_resource_group.hideNsneak.location}"
  resource_group_name          = "${azurerm_resource_group.hideNsneak.name}"
  public_ip_address_allocation = "static"
}

resource "azurerm_virtual_network" "hideNsneak" {
  name                = "hideNsneakNet${random_string.net_name.result}"
  count               = "${var.azure_instance_count > 0 ? 1 : 0}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.hideNsneak.location}"
  resource_group_name = "${azurerm_resource_group.hideNsneak.name}"
}

resource "azurerm_subnet" "hideNsneak" {
  name                 = "hideNsneak${random_string.net_name.result}"
  count                = 0
  resource_group_name  = "${azurerm_resource_group.hideNsneak.name}"
  virtual_network_name = "${azurerm_virtual_network.hideNsneak.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "hideNsneak" {
  count               = "${var.azure_instance_count > 0 ? 1 : 0}"
  location            = "${azurerm_resource_group.hideNsneak.location}"
  resource_group_name = "${azurerm_resource_group.hideNsneak.name}"

  ip_configuration {
    name                          = "hideNsneakSubnet${random_string.net_name.result}"
    subnet_id                     = "${azurerm_subnet.hideNsneak.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_managed_disk" "hideNsneak" {
  name                 = "datadisk_existing"
  count                = "${var.azure_instance_count > 0 ? 1 : 0}"
  location             = "${azurerm_resource_group.hideNsneak.location}"
  resource_group_name  = "${azurerm_resource_group.hideNsneak.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

#azurerm_virtual_machine.hideNsneak.*.ip_adress
resource "azurerm_virtual_machine" "hideNsneak" {
  name                  = "hideNsneak"
  count                 = "${var.azure_instance_count}"
  location              = "${azurerm_resource_group.hideNsneak.location}"
  resource_group_name   = "${azurerm_resource_group.hideNsneak.name}"
  network_interface_ids = ["${azurerm_network_interface.hideNsneak.id}"]
  vm_size               = "${var.azure_vm_size}"

  # Uncomment this line to delete the OS disk automatically when deleting the VM
  delete_os_disk_on_termination = true

  # Uncomment this line to delete the data disks automatically when deleting the VM
  delete_data_disks_on_termination = true

  //TODO: Make dynamic referece to storage image so use can specify
  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk1"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${var.azure_admin_hostname}${azurerm_virtual_machine.hideNsneak.count}"
    admin_username = "${var.azure_default_username}"
  }

  os_profile_linux_config {
    disable_password_authentication = true

    ssh_keys {
      path     = "/home/${var.azure_default_username}/.ssh/authorized_keys"
      key_data = "${file(var.azure_public_key_file)}"
    }
  }

  tags {
    environment = "${var.azure_environment}"
  }

  //TODO: Figure out how to do dynamic provisioner for azure
  # provisioner "local-exec" {
  #   command = "sleep 120; ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -u ${var.azure_default_username} --private-key ${var.azure_private_key_file} -i '${self.},' ../ansible/setup.yml"
  # }
}
