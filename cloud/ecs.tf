resource "aws_ecs_cluster" "platform" {
  name = replace(local.name, ".", "-")
}

resource "aws_ecs_task_definition" "grafana" {
  family                   = "grafana"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 512
  memory                   = 1024
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "grafana"
      image     = "grafana/grafana:9.4.7"
      cpu       = 512
      memory    = 1024
      essential = true
      environment = [
        {
          name  = "GF_AUTH_DISABLE_LOGIN_FORM"
          value = "true"
        },
        {
          name  = "GF_AUTH_DISABLE_SIGNOUT_MENU"
          value = "true"
        },
        {
          name  = "GF_AUTH_ANONYMOUS_ENABLED"
          value = "true"
        },
        {
          name  = "GF_AUTH_ANONYMOUS_ORG_ROLE"
          value = "Viewer"
        },
      ],
      portMappings = [
        {
          containerPort = 3000
          hostPort      = 3000
          protocol      = "tcp"
        }
      ]
    }
  ])
}

resource "aws_ecs_service" "grafana" {
  name                               = "grafana"
  cluster                            = aws_ecs_cluster.platform.id
  task_definition                    = aws_ecs_task_definition.grafana.arn
  desired_count                      = 1
  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200
  health_check_grace_period_seconds  = 15
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"

  network_configuration {
    security_groups  = [aws_security_group.grafana.id]
    subnets          = [for subnet in aws_subnet.public : subnet.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.platform.arn
    container_name   = "grafana"
    container_port   = 3000
  }
}
