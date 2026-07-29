#!/usr/bin/env bash

set -euo pipefail

required_variables=(
  AWS_REGION
  ECS_CLUSTER
  MIGRATION_TASK_DEFINITION
  PRIVATE_SUBNET_IDS
  ECS_SECURITY_GROUP_ID
)

for variable_name in "${required_variables[@]}"; do
  if [[ -z "${!variable_name:-}" ]]; then
    echo "missing required environment variable: ${variable_name}" >&2
    exit 1
  fi
done

network_configuration="awsvpcConfiguration={subnets=[${PRIVATE_SUBNET_IDS}],securityGroups=[${ECS_SECURITY_GROUP_ID}],assignPublicIp=DISABLED}"

echo "Starting migration task ${MIGRATION_TASK_DEFINITION} in cluster ${ECS_CLUSTER}"
task_arn="$(
  aws ecs run-task \
    --region "${AWS_REGION}" \
    --cluster "${ECS_CLUSTER}" \
    --task-definition "${MIGRATION_TASK_DEFINITION}" \
    --launch-type FARGATE \
    --network-configuration "${network_configuration}" \
    --query 'tasks[0].taskArn' \
    --output text
)"

if [[ -z "${task_arn}" || "${task_arn}" == "None" ]]; then
  echo "ECS did not start the migration task" >&2
  exit 1
fi

echo "Waiting for migration task ${task_arn}"
aws ecs wait tasks-stopped \
  --region "${AWS_REGION}" \
  --cluster "${ECS_CLUSTER}" \
  --tasks "${task_arn}"

exit_code="$(
  aws ecs describe-tasks \
    --region "${AWS_REGION}" \
    --cluster "${ECS_CLUSTER}" \
    --tasks "${task_arn}" \
    --query "tasks[0].containers[?name=='migration'].exitCode | [0]" \
    --output text
)"

if [[ "${exit_code}" != "0" ]]; then
  stopped_reason="$(
    aws ecs describe-tasks \
      --region "${AWS_REGION}" \
      --cluster "${ECS_CLUSTER}" \
      --tasks "${task_arn}" \
      --query 'tasks[0].stoppedReason' \
      --output text
  )"
  container_reason="$(
    aws ecs describe-tasks \
      --region "${AWS_REGION}" \
      --cluster "${ECS_CLUSTER}" \
      --tasks "${task_arn}" \
      --query "tasks[0].containers[?name=='migration'].reason | [0]" \
      --output text
  )"
  echo "Migration failed with exit code ${exit_code}; task reason: ${stopped_reason}; container reason: ${container_reason}" >&2
  exit 1
fi

echo "Migration completed successfully; service rollout may continue"
