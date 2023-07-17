provider "google" {
  credentials = file("credentials.json")
  project     = "blueteam-vargroup"
  region      = "europe-west8"
}

resource "google_compute_network" "bt_gke" {
  name                    = "bt-gke"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "node" {
  name                    = "node"
  ip_cidr_range           = "10.19.0.0/24"
  stack_type              = "IPV4_ONLY"
  network                 = google_compute_network.bt_gke.self_link
  region                  = "europe-west8"
  private_ip_google_access = true
}

resource "google_compute_address" "mosquitto_ip_statico" {
  name    = "mosquitto-ip-statico"
  region  = "europe-west8"
}

resource "google_compute_disk" "mosquitto" {
  name  = "mosquitto"
  type  = "pd-balanced"
  size  = 10
  zone  = "europe-west8-a"
}

resource "google_container_cluster" "blueteam" {
  name                     = "blueteam"
  location                 = "europe-west8"
  remove_default_node_pool = true

  network = google_compute_network.bt_gke.self_link
  subnetwork = google_compute_subnetwork.node.self_link

  node_pool {
    name               = "default-pool"
    initial_node_count = 1
    version            = "1.26.5-gke.1200"
  }
}

resource "google_bigquery_dataset" "blueteam_mqtt" {
  dataset_id = "blueteam_mqtt"
  location   = "EU"
}

resource "google_bigquery_table" "eventi" {
  dataset_id = google_bigquery_dataset.blueteam_mqtt.dataset_id
  table_id   = "eventi"

    schema = <<EOF
[
  {
    "name": "uuid",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "impianto",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "codice",
    "type": "INT64",
    "mode": "REQUIRED"
  },
  {
    "name": "descrizione",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "valore",
    "type": "STRING",
    "mode": "REQUIRED"
  },
  {
    "name": "data",
    "type": "DATETIME",
    "mode": "REQUIRED"
  }
]
EOF
}

