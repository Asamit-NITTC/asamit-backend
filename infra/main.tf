variable CLIENT_ID {}
variable TEST_SUB {}
variable CLOUD_SQL_USER_NAME{}
variable CLOUD_SQL_PASSWORD{}
variable CLOUD_SQL_IP{}
variable CLOUD_SQL_PORT{}



provider "google" {
    project         = "asami-1e0c2"
    region          = "asia-northeast1"
    credentials     = file("../key.json")
    request_timeout = "60s"
}

resource "google_cloud_run_v2_service" "default" {
  name     = "asamit-cloudrun-terraform"
  location = "asia-northeast1"
  ingress = "INGRESS_TRAFFIC_ALL"


  template{

    containers {
      image = "asia-northeast1-docker.pkg.dev/asami-1e0c2/terraform-test/asamit:latest"
      env {
        name = "CLIENT_ID"
        value = "${var.CLIENT_ID}"
      }
      env {
        name = "TEST_SUB"
        value = "${var.TEST_SUB}"
      }
      env {
        name = "CLOUD_SQL_USER_NAME"
        value = "${var.CLOUD_SQL_USER_NAME}"
      }
      env {
        name = "CLOUD_SQL_PASSWORD"
        value = "${var.CLOUD_SQL_PASSWORD}"
      }
      env {
        name = "CLOUD_SQL_IP"
        value = "${var.CLOUD_SQL_IP}"
      }
      env {
        name = "CLOUD_SQL_PORT"
        value = "${var.CLOUD_SQL_PORT}"
      }
    }

  }
}

resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloud_run_v2_service.default.location
  service  = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
