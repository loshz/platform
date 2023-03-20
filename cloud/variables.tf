locals {
  name = "platform.loshz.com"
  tags = {
    "terraform.io/managed"       = "true"
    "platform.loshz.com/version" = var.platform_version
  }
}

variable "region" {
  type        = string
  description = "AWS region in which all resources will be created"
}

variable "platform_version" {
  type        = string
  description = "Version representing the current resource state"
  default     = "v0.1.0"
}

variable "vpc_cidr" {
  type        = string
  description = "CIDR of the VPC"
  default     = "10.0.0.0/16"

  validation {
    condition     = can(cidrnetmask(var.vpc_cidr))
    error_message = "Must be a valid IPv4 CIDR block address."
  }
}
