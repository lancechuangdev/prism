variable "aws_region" {
  description = "AWS region for all resources."
  type        = string
}

variable "environment" {
  description = "Deployment environment name."
  type        = string
  default     = "production"

  validation {
    condition     = length(trimspace(var.environment)) > 0
    error_message = "environment must not be empty."
  }
}

variable "vpc_cidr" {
  description = "CIDR assigned to the VPC."
  type        = string
  default     = "10.42.0.0/16"
}

variable "image_uri" {
  description = "Immutable backend container image URI, preferably pinned by digest."
  type        = string
}

variable "domain_name" {
  description = "API DNS name, such as api.prism.example."
  type        = string

  validation {
    condition     = length(trimspace(var.domain_name)) > 0
    error_message = "domain_name must not be empty."
  }
}

variable "route53_zone_id" {
  description = "Route 53 hosted-zone ID containing domain_name."
  type        = string
}

variable "chain_id" {
  description = "Expected EVM chain ID."
  type        = string
}

variable "chain_rpc_url_secret_arn" {
  description = "Secrets Manager ARN whose value is the EVM JSON-RPC URL."
  type        = string
}

variable "pool_address" {
  description = "PrismPool address from the verified deployment manifest."
  type        = string

  validation {
    condition     = can(regex("^0x[0-9a-fA-F]{40}$", var.pool_address))
    error_message = "pool_address must be a 20-byte hexadecimal EVM address."
  }
}

variable "multisig_address" {
  description = "ThresholdMultiSig address from the verified deployment manifest."
  type        = string

  validation {
    condition     = can(regex("^0x[0-9a-fA-F]{40}$", var.multisig_address))
    error_message = "multisig_address must be a 20-byte hexadecimal EVM address."
  }
}

variable "price_provider_url" {
  description = "HTTPS quote-provider endpoint."
  type        = string

  validation {
    condition     = startswith(lower(var.price_provider_url), "https://")
    error_message = "price_provider_url must use HTTPS."
  }
}

variable "price_provider_token_secret_arn" {
  description = "Secrets Manager ARN whose value is the quote-provider bearer token."
  type        = string
}

variable "cognito_region" {
  type    = string
  default = ""
}

variable "cognito_user_pool_id" {
  type    = string
  default = ""
}

variable "cognito_client_id" {
  type    = string
  default = ""
}

variable "db_name" {
  type    = string
  default = "prism"
}

variable "db_username" {
  type    = string
  default = "prism"
}

variable "redis_auth_token_secret_arn" {
  description = "Secrets Manager ARN whose value is the Redis AUTH token."
  type        = string
}

variable "api_desired_count" {
  type    = number
  default = 2

  validation {
    condition     = var.api_desired_count >= 2
    error_message = "api_desired_count must be at least 2 for multi-AZ availability."
  }
}

variable "api_min_capacity" {
  type    = number
  default = 2
}

variable "api_max_capacity" {
  type    = number
  default = 6
}
