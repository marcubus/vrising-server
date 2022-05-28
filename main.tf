terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.3.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.1.0"
    }
  }
}

resource "random_pet" "vrising" {
  keepers = {
    vrising = var.world_name
  }
}

provider "aws" {
  region = var.region
}

resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}

locals {
  world  = "vrising-${random_pet.vrising.id}"
  vpc_id = var.vpc_id == "" ? aws_default_vpc.default.id : var.vpc_id
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [local.vpc_id]
  }
}

resource "aws_s3_bucket" "backups" {
  bucket_prefix = "vrising-backup-${lower(var.world_name)}"
}

resource "aws_s3_bucket_acl" "backups" {
  bucket = aws_s3_bucket.backups.id
  acl    = "private"
}

data aws_caller_identity account {
}
