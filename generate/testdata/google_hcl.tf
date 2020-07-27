resource "google_compute_instance" "inframap-tmp" {
  name         = "inframap-tmp"
  machine_type = "f1-micro"
  zone         = "us-east1-b"

  tags = ["front", "ssh"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  scheduling {
    preemptible       = true
    automatic_restart = false
  }

  network_interface {
    network = google_compute_network.vpc_network.name

    access_config {
      // Ephemeral IP
    }
  }

}

resource "google_compute_instance" "inframap-tmp-two" {
  name         = "inframap-tmp-two"
  machine_type = "f1-micro"
  zone         = "us-east1-b"

  tags = ["ssh"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  scheduling {
    preemptible       = true
    automatic_restart = false
  }

  network_interface {
    network = google_compute_network.vpc_network.name

    access_config {
      // Ephemeral IP
    }
  }

}

resource "google_compute_network" "vpc_network" {
  name                    = "inframap-network-tmp"
  auto_create_subnetworks = true
}

resource "google_compute_firewall" "allow-http" {
  network     = google_compute_network.vpc_network.name
  description = "allow incoming http traffic"
  name        = "allow-http"

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["front"]
}

resource "google_compute_firewall" "allow-ssh" {
  network     = google_compute_network.vpc_network.name
  description = "allow SSH between instances"
  name        = "allow-ssh"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  direction   = "INGRESS"
  source_tags = ["ssh"]
  target_tags = ["front"]
}
