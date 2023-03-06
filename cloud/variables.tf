locals {
  tags = {
    "terraform.io/managed"       = "true"
    "platform.loshz.com/version" = var.version
  }
}

variable "region" {
  type        = string
  description = "AWS region in which all resources will be created"
}

variable "version" {
  type        = string
  description = "Version representing the current resource state"
  default     = "v0.1.0"
}
