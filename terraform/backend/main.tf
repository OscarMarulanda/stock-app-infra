resource "aws_instance" "go_api" {
  ami                    = "ami-0d406e26e5ad4de53"  # Amazon Linux 2
  instance_type          = "t3.micro"
  key_name               = var.key_name
  associate_public_ip_address = true

  user_data = file("${path.module}/install.sh")

  tags = {
    Name = "GoBackend"
  }
}