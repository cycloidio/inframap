data "flexibleengine_images_image_v2" "ubuntu" {
  name        = "Ubuntu 18.04"
  most_recent = true
}

resource "flexibleengine_compute_instance_v2" "instance_one" {
  name      = "instance-one"
  image_id  = "${data.flexibleengine_images_image_v2.ubuntu.id}"
  flavor_id = "s3.small.1"

  network {
    port = "${flexibleengine_networking_port_v2.port_instance_one.id}"
  }
}

resource "flexibleengine_compute_instance_v2" "instance_two" {
  name      = "instance-two"
  image_id  = "${data.flexibleengine_images_image_v2.ubuntu.id}"
  flavor_id = "s3.small.1"

  network {
    port = "${flexibleengine_networking_port_v2.port_instance_two.id}"
  }
}

resource "flexibleengine_networking_port_v2" "port_instance_one" {
  name = "port_instance_one"
  security_group_ids = [
    "${flexibleengine_networking_secgroup_v2.secgroup_instance_one.id}",
  ]
}

resource "flexibleengine_networking_port_v2" "port_instance_two" {
  name = "port_instance_two"
  security_group_ids = [
    "${flexibleengine_networking_secgroup_v2.secgroup_instance_two.id}",
  ]
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_instance_one" {
  name        = "secgroup_instance_one"
  description = "security group for the instance one"
}

resource "flexibleengine_networking_secgroup_v2" "secgroup_instance_two" {
  name        = "secgroup_instance_two"
  description = "security group for the instance two"
}

resource "flexibleengine_networking_secgroup_rule_v2" "ssh_two_to_one" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  security_group_id = "${flexibleengine_networking_secgroup_v2.secgroup_instance_one.id}"
  remote_group_id   = "${flexibleengine_networking_secgroup_v2.secgroup_instance_two.id}"
}
