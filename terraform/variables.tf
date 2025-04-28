variable "aws_region" {
  default = "us-east-2"
}

variable "key_name" {
  description = "EC2 Key pair name"
  type        = string
}

variable "alpha_vantage_api_key" {
  description = "Alpha Vantage API key"
  type        = string
  sensitive   = true
}

variable "cockroachdb_dsn" {
  description = "CockroachDB connection string"
  type        = string
  sensitive   = true
}