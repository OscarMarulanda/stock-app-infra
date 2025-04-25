module "ssm" {
  source = "./ssm"
  alpha_vantage_api_key = var.alpha_vantage_api_key
}

module "frontend" {
  source = "./frontend"
}

module "backend" {
  source = "./backend"
  key_name = var.key_name
}