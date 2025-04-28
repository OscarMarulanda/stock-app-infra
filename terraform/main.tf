

module "frontend" {
  source = "./frontend"
}

module "backend" {
  source              = "./backend"
  alpha_vantage_api_key = var.alpha_vantage_api_key
  key_name           = var.key_name
  cockroachdb_dsn    = var.cockroachdb_dsn
}


output "go_api_public_ip" {
  description = "The public IP of the EC2 instance running the Go API"
  value       = module.backend.go_api_public_ip
}