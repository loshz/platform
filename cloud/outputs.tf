output "public_load_balancer" {
  value = aws_lb.platform.dns_name
}
