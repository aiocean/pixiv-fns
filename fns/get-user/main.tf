variable "name" {}

variable "gcp_project_id" {}
variable "gcp_region" {}
variable "service_account_email" {}

data "archive_file" "source" {
  type        = "zip"
  source_dir  = abspath("${path.module}/")
  output_path = "/tmp/function.zip"
}

resource "google_storage_bucket" "bucket" {
  name     = "${var.gcp_project_id}-${lower(var.name)}-gcf-source"
  location = var.gcp_region
}

resource "google_storage_bucket_object" "zip" {
  name   = "${data.archive_file.source.output_md5}.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.source.output_path
}

resource "google_cloudfunctions_function" "function" {
  available_memory_mb = "128"
  entry_point         = "Handler"
  ingress_settings    = "ALLOW_ALL"

  name                  = var.name
  project               = var.gcp_project_id
  region                = var.gcp_region
  runtime               = "go116"
  service_account_email = var.service_account_email
  timeout               = 20
  trigger_http          = true
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = "${data.archive_file.source.output_md5}.zip"
}

resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.function.project
  region         = google_cloudfunctions_function.function.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

output "function_url" {
  value = google_cloudfunctions_function.function.https_trigger_url
}
