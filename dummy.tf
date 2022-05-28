resource "aws_ecs_cluster" "dummy" {
  name = "dummy-${local.world}"
}

resource "aws_ecs_service" "dummy" {
  name            = "dummy-${local.world}"
  cluster         = aws_ecs_cluster.dummy.name
  task_definition = aws_ecs_task_definition.dummy.arn
  desired_count   = 0

  deployment_maximum_percent         = 100
  deployment_minimum_healthy_percent = 0

  load_balancer {
    target_group_arn = aws_lb_target_group.fiveseven.arn
    container_name   = "dummy"
    container_port   = 9876
  }
}

resource "aws_ecs_task_definition" "dummy" {
  family       = "dummy-${local.world}"
  network_mode = "bridge"
  cpu          = 256
  memory       = 256

  container_definitions = jsonencode([
    {
      "name"      = "dummy",
      "image"     = "marcubus/vrising-server-dummy:latest",
      "cpu"       = 256,
      "memory"    = 256,
      "essential" = true,
      "environment" = [
        {
          "name"  = "TARGET_SERVICE",
          "value" = local.world
        },
        {
          "name"  = "DUMMY_SERVICE",
          "value" = "dummy-${local.world}"
        },
        {
          "name"  = "TARGET_CLUSTER",
          "value" = local.world
        },
        {
          "name"  = "DUMMY_CLUSTER",
          "value" = "dummy-${local.world}"
        },
        {
          "name"  = "TARGET_ASG",
          "value" = local.world
        },
        {
          "name"  = "DUMMY_ASG",
          "value" = "dummy-${local.world}"
        },
        {
          "name"  = "AWS_REGION",
          "value" = var.region
        }
      ],
      "portMappings" = [
        {
          "containerPort" = 9876,
          "hostPort"      = 9876,
          "protocol"      = "udp"
        }
      ]
    }
  ])
}
