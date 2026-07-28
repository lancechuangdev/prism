# Prism AWS infrastructure

This Terraform stack creates the production AWS baseline: a two-AZ VPC, public HTTPS Application Load Balancer, private ECS/Fargate API and scheduler, private Multi-AZ RDS MySQL, encrypted ElastiCache Redis, Route 53 DNS, ACM TLS, CloudWatch logs, IAM task roles, API autoscaling, and a one-shot migration task definition.

Copy `terraform.tfvars.example` to an untracked `terraform.tfvars`, use an
immutable container digest, and run:

```bash
cp backend.hcl.example backend.hcl
# Fill in the pre-created, versioned S3 state bucket and DynamoDB lock table.
terraform init -backend-config=backend.hcl
terraform plan -out=tfplan
terraform apply tfplan
```

Run migrations before updating the services:

```bash
aws ecs run-task \
  --cluster "$(terraform output -raw ecs_cluster_name)" \
  --task-definition "$(terraform output -raw migration_task_definition_arn)" \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[$(terraform output -json private_subnet_ids | jq -r 'join(\",\")')],securityGroups=[$(terraform output -raw ecs_security_group_id)],assignPublicIp=DISABLED}"
```

Wait for the migration task to exit successfully before deploying a new API or
scheduler revision. Terraform state contains sensitive input values in this
initial stack, so the stack requires an encrypted, access-controlled S3 remote
backend rather than silently creating local production state. The state bucket
and lock table are bootstrap resources and must exist before `terraform init`.
Moving runtime secrets out of task-definition environment variables and into
AWS Secrets Manager is intentionally tracked as the next production blocker.
