locals {
  metrics_namespace = "Prism/${var.environment}"
}

resource "aws_sns_topic" "alerts" {
  name = "${local.name}-alerts"
}

resource "aws_sns_topic_subscription" "email" {
  topic_arn = aws_sns_topic.alerts.arn
  protocol  = "email"
  endpoint  = var.alarm_email
}

resource "aws_cloudwatch_log_metric_filter" "api_5xx" {
  name           = "${local.name}-api-5xx"
  log_group_name = aws_cloudwatch_log_group.api.name
  pattern        = "{ $.event = \"http_request\" && $.status >= 500 }"

  metric_transformation {
    name      = "API5xxCount"
    namespace = local.metrics_namespace
    value     = "1"
  }
}

resource "aws_cloudwatch_log_metric_filter" "scheduler_failure" {
  name           = "${local.name}-scheduler-failure"
  log_group_name = aws_cloudwatch_log_group.scheduler.name
  pattern        = "{ $.event = \"scheduler_sync_failure\" }"

  metric_transformation {
    name      = "SchedulerFailureCount"
    namespace = local.metrics_namespace
    value     = "1"
  }
}

resource "aws_cloudwatch_log_metric_filter" "scheduler_success" {
  name           = "${local.name}-scheduler-success"
  log_group_name = aws_cloudwatch_log_group.scheduler.name
  pattern        = "{ $.event = \"scheduler_sync_success\" }"

  metric_transformation {
    name      = "SchedulerSuccessCount"
    namespace = local.metrics_namespace
    value     = "1"
  }
}

resource "aws_cloudwatch_log_metric_filter" "provider_failure_scheduler" {
  name           = "${local.name}-provider-failure-scheduler"
  log_group_name = aws_cloudwatch_log_group.scheduler.name
  pattern        = "{ $.event = \"provider_failure\" }"

  metric_transformation {
    name      = "ProviderFailureCount"
    namespace = local.metrics_namespace
    value     = "1"
  }
}

resource "aws_cloudwatch_log_metric_filter" "migration_failure" {
  name           = "${local.name}-migration-failure"
  log_group_name = aws_cloudwatch_log_group.migration.name
  pattern        = "{ $.event = \"migration_failure\" }"

  metric_transformation {
    name      = "MigrationFailureCount"
    namespace = local.metrics_namespace
    value     = "1"
  }
}

resource "aws_cloudwatch_metric_alarm" "api_5xx" {
  alarm_name          = "${local.name}-api-5xx"
  alarm_description   = "The Prism API returned at least one 5xx response."
  namespace           = local.metrics_namespace
  metric_name         = "API5xxCount"
  statistic           = "Sum"
  period              = 300
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanOrEqualToThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "scheduler_failure" {
  alarm_name          = "${local.name}-scheduler-failure"
  alarm_description   = "The scheduler failed to complete a synchronization."
  namespace           = local.metrics_namespace
  metric_name         = "SchedulerFailureCount"
  statistic           = "Sum"
  period              = 300
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanOrEqualToThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "scheduler_lag" {
  alarm_name          = "${local.name}-scheduler-lag"
  alarm_description   = "No successful scheduler synchronization was observed for two consecutive five-minute periods."
  namespace           = local.metrics_namespace
  metric_name         = "SchedulerSuccessCount"
  statistic           = "Sum"
  period              = 300
  evaluation_periods  = 2
  datapoints_to_alarm = 2
  threshold           = 1
  comparison_operator = "LessThanThreshold"
  treat_missing_data  = "breaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "provider_failure" {
  alarm_name          = "${local.name}-provider-failure"
  alarm_description   = "The scheduler reported an upstream chain RPC or ChainlinkOracle price-read failure."
  namespace           = local.metrics_namespace
  metric_name         = "ProviderFailureCount"
  statistic           = "Sum"
  period              = 300
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanOrEqualToThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "migration_failure" {
  alarm_name          = "${local.name}-migration-failure"
  alarm_description   = "The one-shot migration task reported a failure."
  namespace           = local.metrics_namespace
  metric_name         = "MigrationFailureCount"
  statistic           = "Sum"
  period              = 60
  evaluation_periods  = 1
  threshold           = 1
  comparison_operator = "GreaterThanOrEqualToThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]
}

resource "aws_cloudwatch_metric_alarm" "api_cpu" {
  alarm_name          = "${local.name}-api-high-cpu"
  alarm_description   = "Average ECS API service CPU utilization remained above 80% for ten minutes."
  namespace           = "AWS/ECS"
  metric_name         = "CPUUtilization"
  statistic           = "Average"
  period              = 300
  evaluation_periods  = 2
  threshold           = 80
  comparison_operator = "GreaterThanThreshold"
  treat_missing_data  = "notBreaching"
  alarm_actions       = [aws_sns_topic.alerts.arn]
  ok_actions          = [aws_sns_topic.alerts.arn]

  dimensions = {
    ClusterName = aws_ecs_cluster.main.name
    ServiceName = aws_ecs_service.api.name
  }
}
