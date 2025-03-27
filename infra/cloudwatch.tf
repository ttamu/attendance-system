resource "aws_cloudwatch_log_group" "ecs_api_log_group" {
  name              = "/ecs/${var.project_name}-api"
  retention_in_days = 7
}
