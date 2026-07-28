output "api_url" {
  value = "https://${var.domain_name}"
}

output "ecs_cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "migration_task_definition_arn" {
  value = aws_ecs_task_definition.migration.arn
}

output "private_subnet_ids" {
  value = values(aws_subnet.private)[*].id
}

output "ecs_security_group_id" {
  value = aws_security_group.ecs.id
}

