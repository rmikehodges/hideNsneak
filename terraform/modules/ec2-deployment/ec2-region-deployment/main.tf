provider "aws" {}

locals {
  keyCount = "${var.instance_count > 0 ? 1 : 0}"
}

data "aws_subnet_ids" "all" {
  count  = "${var.instance_count > 0 ? 1 : 0}"
  vpc_id = "${data.aws_vpc.default.id}"
}

data "aws_vpc" "default" {
  count   = "${var.instance_count > 0 ? 1 : 0}"
  default = true
}

data "aws_ami" "ubuntu" {
  count       = "${var.instance_count > 0 ? 1 : 0}"
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

data "aws_security_groups" "hideNsneak" {
  count = "${var.instance_count > 0 ? 1 : 0}"

  filter {
    name   = "group-name"
    values = ["hidensneak"]
  }
}

resource "aws_instance" "hideNsneak" {
  ami = "${data.aws_ami.ubuntu.id}"

  instance_type = "${var.aws_instance_type == "" ? "t2.micro" :  var.aws_instance_type}"

  count     = "${var.instance_count}"
  subnet_id = "${element(data.aws_subnet_ids.all.ids, 0)}"

  key_name = "${var.aws_keypair_name}"

  vpc_security_group_ids = ["${data.aws_security_groups.hideNsneak.ids}"]

  tags {
    Name = "hidensneak"
  }

  depends_on = ["data.aws_ami.ubuntu", "data.aws_vpc.default", "data.aws_subnet_ids.all"]
}

resource "ansible_host" "hideNsneak" {
  count = "${var.instance_count}"

  inventory_hostname = "${aws_instance.hideNsneak.*.public_ip[count.index]}"
  groups             = "${var.ansible_groups}"

  vars {
    ansible_user                 = "${var.ec2_default_user}"
    ansible_connection           = "ssh"
    ansible_ssh_private_key_file = "${var.aws_private_key_file}"
    ansible_ssh_common_args      = "-o StrictHostKeyChecking=no"
  }
}
