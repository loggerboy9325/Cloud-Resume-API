terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  backend "remote" {
    organization = "gwettlaufertest"

    workspaces {
      name = "cloud-resume-api"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}
