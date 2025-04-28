variable "key_name" {
  type        = string
  description = "The EC2 key pair name for SSH access"
}

variable "alpha_vantage_api_key" {
  description = "API key for Alpha Vantage"
  type        = string
  sensitive   = true
}

variable "cockroachdb_dsn" {
  description = "CockroachDB connection string"
  type        = string
  sensitive   = true
}