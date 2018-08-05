# resource "ansible_host" "hideNsneak" {
#   count = "${var.do_count}"

#   //Element
#   inventory_hostname = "${digitalocean_droplet.hideNsneak.*.ipv4_address[count.index]}"

#   //Possibly add groups in the future
#   # groups             = "${var.ansible_groups}"

#   vars {
#     ansible_user                 = "${var.do_default_user}"
#     ansible_connection           = "ssh"
#     ansible_ssh_private_key_file = "${var.do_private_key}"
#     ansible_ssh_common_args      = "-o StrictHostKeyChecking=no"
#   }
#   depends_on = ["digitalocean_droplet.hideNsneak"]
# }

resource "digitalocean_droplet" "hideNsneak" {
  image  = "${var.do_image == "" ? "ubuntu-16-04-x64" : var.do_image}"
  name   = "${var.name}"
  region = "${var.do_region}"
  size   = "${var.do_size == "" ? "512mb" : var.do_size}"
  count  = "${var.do_count}"

  ssh_keys = [
    "${var.do_ssh_fingerprint}",
  ]

  //Uncomment this for ansible
  # provisioner "local-exec" {
  #   command = "sleep 120; ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -u ${var.do_default_user} --private-key ${var.do_private_key} -i '${self.ipv4_address},' ../ansible/setup.yml"
  # }
}

//Uncomment this for default firewall rules
# resource "digitalocean_firewall" "hideNsneak" {
#   name        = "hidensneak-test"
#   droplet_ids = ["${digitalocean_droplet.hideNsneak.*.id}"]
#   count       = "${digitalocean_droplet.hideNsneak.count > 0 ? 1 : 0}"


#   inbound_rule = [
#     {
#       protocol         = "tcp"
#       port_range       = "22"
#       source_addresses = ["0.0.0.0/0"]
#     },
#   ]


#   outbound_rule = [
#     {
#       protocol              = "tcp"
#       port_range            = "1-65535"
#       destination_addresses = ["0.0.0.0/0", "::/0"]
#     },
#     {
#       protocol              = "udp"
#       port_range            = "1-65535"
#       destination_addresses = ["0.0.0.0/0", "::/0"]
#     },
#   ]
# }

