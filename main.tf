terraform {
  cloud {
    organization = "aiocean"
    workspaces {
      name = "pixiv"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  required_version = ">= 0.14.0"
}

locals {
  functions = {
    "get-artwork" = {
      name = "get-artwork"
    }
  }
}

provider "aws" {
  region     = var.aws_region
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key

  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_requesting_account_id  = true
}

module "lambda_function" {
  for_each = local.functions

  source        = "terraform-aws-modules/lambda/aws"
  architectures = ["arm64"]
  function_name = "${terraform.workspace}-${var.env}-${each.value.name}"
  runtime       = "provided.al2"
  handler       = "./bootstrap"
  source_path   = "${path.module}/fns/${each.value.name}/.bin/bootstrap"
}
