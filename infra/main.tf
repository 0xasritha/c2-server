# umassctf25@gmail.com/project-asritha-kat
# authenticate w/ `gcloud auth application-default login`

provider "google" {
  project = "project-asritha-kat" # project name: `564-project-asritha-kat`
  region = "us-central1"
  zone = "us-central1-c"
}

# C2 Server - disguised as "Name Server"
resource "google_compute_instance" "c2_server" {
  name         = "c2-server"
  machine_type = "e2-medium"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network = "default"

    access_config {
      nat_ip = google_compute_address.c2_server_static.address
    }
  }

  tags = ["http-server", "https-server"]
}

# Static External IP
resource "google_compute_address" "c2_server_static" {
  name = "c2-server-static"
}

# Proxy - #TODO

# Firewall Rules
# TODO: expose port 53 ? -> figure out setup behind proxy
resource "google_compute_firewall" "allow_http" {
  name    = "default-allow-http"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["http-server"]
  direction     = "INGRESS"
  description   = "Allow HTTP traffic"
}

resource "google_compute_firewall" "allow_https" {
  name    = "default-allow-https"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["https-server"]
  direction     = "INGRESS"
  description   = "Allow HTTPS traffic"
}

resource "google_compute_firewall" "allow_dns" {
  name    = "default-allow-dns"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["53"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["dns-server"]
  direction     = "INGRESS"
  description   = "Allow DNS traffic"
}

# Storage Bucket to host implant
resource "google_storage_bucket" "implant" {
  name = "564-proj-host-implant"
  location = "US"
}

resource "google_storage_bucket_iam_member" "make_public" {
  bucket   = google_storage_bucket.implant.name
  role     = "roles/storage.objectViewer"
  member   = "allUsers"
}

# Private Container Registry
#TODO