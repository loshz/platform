# Create a public load balancer to serve the root platform domain.
resource "aws_lb" "platform" {
  name               = "platform"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.load_balancer.id]
  subnets            = [for subnet in aws_subnet.public : subnet.id]

  enable_deletion_protection = false
}

# Create a target group for the Grafana containers routed to from
# the load baancer.
resource "aws_lb_target_group" "platform" {
  name        = "platform"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = aws_vpc.platform.id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    path                = "/api/health"
    port                = 3000
    protocol            = "HTTP"
    timeout             = 5
    healthy_threshold   = 5
    unhealthy_threshold = 2
    matcher             = 200
  }
}

resource "aws_lb_listener" "platform_http" {
  load_balancer_arn = aws_lb.platform.id
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"

    redirect {
      port        = 443
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_alb_listener" "platform_https" {
  load_balancer_arn = aws_lb.platform.id
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS-1-1-2017-01"
  certificate_arn   = var.platform_cert_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.platform.arn
  }
}
