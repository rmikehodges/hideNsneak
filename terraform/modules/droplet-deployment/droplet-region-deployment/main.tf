resource "digitalocean_droplet" "hideNsneak" {
  image  = "${var.do_image == "" ? "ubuntu-16-04-x64" : var.do_image}"
  name   = "${var.name}"
  region = "${var.do_region}"
  size   = "${var.do_size == "" ? "512mb" : var.do_size}"
  count  = "${var.do_count}"

  ssh_keys = [
    "${var.do_ssh_fingerprint}",
  ]
}