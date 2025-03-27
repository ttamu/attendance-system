output "vpc_id" {
  description = "ID of the created VPC"
  value       = aws_vpc.main.id
}

# output "ecs_cluster_id" {
#   description = "ID of the created ECS Cluster"
#   value       = aws_ecs_cluster.main.id
# }
#
# output "alb_dns_name" {
#   description = "DNS name of the Application Load Balancer"
#   value       = aws_lb.app_alb.dns_name
# }
#
# output "rds_endpoint" {
#   description = "Endpoint of the RDS instance"
#   value       = aws_db_instance.db.address
# }
