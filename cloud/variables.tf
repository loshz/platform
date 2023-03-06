locals {
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
