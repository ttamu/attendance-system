resource "aws_ecr_repository" "api_repo" {
  name = "${var.project_name}-api"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name    = "${var.project_name}-api"
    Project = var.project_name
  }
}