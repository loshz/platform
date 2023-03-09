terraform {
  required_version = "~> 1.3"

  backend "s3" {
    key     = "cloud.tfstate"
    encrypt = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.56.0"
    }
  }
}

provider "aws" {
  region = var.region
}
