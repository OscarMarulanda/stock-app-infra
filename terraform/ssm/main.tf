resource "aws_ssm_parameter" "alpha_vantage" {
  name  = "/myapp/alpha_vantage_api_key"
  type  = "SecureString"
  value = var.alpha_vantage_api_key
}