resource "aws_instance" "go_api" {
  ami                    = "ami-0d406e26e5ad4de53"  # Amazon Linux 2
  instance_type          = "t3.micro"
  key_name               = var.key_name
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ec2_instance_profile.name


  user_data = templatefile("${path.module}/install.sh", {
    api_key_ssm_param = aws_ssm_parameter.alpha_vantage_api_key.name
  })

  tags = {
    Name = "GoBackend2"
  }

  depends_on = [aws_ssm_parameter.alpha_vantage_api_key]
}

resource "aws_ssm_parameter" "alpha_vantage_api_key" {
  name        = "/app/alpha_vantage_api_key"
  description = "Alpha Vantage API Key"
  type        = "SecureString"
  value       = var.alpha_vantage_api_key
}

output "go_api_public_ip" {
  description = "The public IP of the EC2 instance running the Go API"
  value       = aws_instance.go_api.public_ip
}