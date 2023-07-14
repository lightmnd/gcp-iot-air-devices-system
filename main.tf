# main.tf

provider "google" {
  credentials = file("path/to/google_credentials.json")
  project     = "gcp-project-id"
  region      = "us-central1"
  zone      = "us-central1-a"
}

resource "google_container_cluster" "gke_cluster" {
  name               = "iot-system-cluster"
  location           = "us-central1"
  initial_node_count = 3
  node_config {
    machine_type = "e2-medium"
  }
}

resource "google_container_node_pool" "gke_node_pool" {
  name       = "iot-system-node-pool"
  cluster    = google_container_cluster.gke_cluster.name
  node_count = 3
  autoscaling {
    min_node_count = 3
    max_node_count = 5
  }
  management {
    auto_repair  = true
    auto_upgrade = true
  }
  node_config {
    machine_type = "e2-medium"
  }
}

resource "helm_release" "mosquitto" {
  name       = "mosquitto"
  repository = "https://charts.helm.sh/stable"
  chart      = "mosquitto"
  namespace  = "default"
  depends_on = [google_container_node_pool.gke_node_pool]
}
