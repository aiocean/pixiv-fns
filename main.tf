terraform {
  cloud {
    organization = "aiocean"
    workspaces {
      name = "pixiv-prod"
    }
  }

  required_version = ">= 0.14.0"
}

variable "gcp_project_id" {}
variable "gcp_region" {}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

provider "google-beta" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

resource "google_service_account" "function-sa" {
  account_id   = "function-${lower(terraform.workspace)}"
  display_name = "function-${lower(terraform.workspace)}"
  project      = var.gcp_project_id
}

module "getUser" {
  source = "./fns/get-user"

  gcp_project_id        = var.gcp_project_id
  gcp_region            = var.gcp_region
  name                  = "${terraform.workspace}-getUser"
  service_account_email = google_service_account.function-sa.email
}

module "getArtwork" {
  source = "./fns/get-artwork"

  gcp_project_id        = var.gcp_project_id
  gcp_region            = var.gcp_region
  name                  = "${terraform.workspace}-getArtwork"
  service_account_email = google_service_account.function-sa.email
}

output "getUser_url" {
  value = module.getUser.function_url
}

output "getArtwork_url" {
  value = module.getArtwork.function_url
}
