
# to init this backend use:
# terraform init -backend-config=$HOME/liftplan-lambda-backend.hcl

terraform {
  backend "remote" {}
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.51.0"
    }
  }
  required_version = ">= 1.0"
}

provider "aws" {
  region = var.aws_region
}
