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

variable "oracle_address" {
  description = "ChainlinkOracle address from the verified deployment manifest."
  type        = string

  validation {
    condition     = can(regex("^0x[0-9a-fA-F]{40}$", var.oracle_address))
    error_message = "oracle_address must be a 20-byte hexadecimal EVM address."
  }
}

variable "price_symbol" {
  description = "Token symbol refreshed by the scheduler."
  type        = string
  default     = "WETH"
}

variable "price_token_addresses" {
  description = "Map of symbols to ERC-20 addresses configured in ChainlinkOracle."
  type        = map(string)

  validation {
    condition = length(var.price_token_addresses) > 0 && alltrue([
      for symbol, address in var.price_token_addresses :
      length(trimspace(symbol)) > 0 && can(regex("^0x[0-9a-fA-F]{40}$", address))
    ])
    error_message = "price_token_addresses must map non-empty symbols to 20-byte hexadecimal EVM addresses."
  }
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

variable "login_rate_limit" {
  description = "Approximate login requests allowed per source IP in AWS WAF's five-minute window."
  type        = number
  default     = 20

  validation {
    condition     = var.login_rate_limit >= 10
    error_message = "login_rate_limit must be at least AWS WAF's minimum of 10."
  }
}

variable "proposal_rate_limit" {
  description = "Approximate proposal requests allowed per source IP in AWS WAF's five-minute window."
  type        = number
  default     = 100

  validation {
    condition     = var.proposal_rate_limit >= 10
    error_message = "proposal_rate_limit must be at least AWS WAF's minimum of 10."
  }
}

variable "alarm_email" {
  description = "Email address subscribed to production CloudWatch alarms. AWS sends a confirmation link before alerts are delivered."
  type        = string

  validation {
    condition     = can(regex("^[^@[:space:]]+@[^@[:space:]]+\\.[^@[:space:]]+$", var.alarm_email))
    error_message = "alarm_email must be a valid email address."
  }
}
