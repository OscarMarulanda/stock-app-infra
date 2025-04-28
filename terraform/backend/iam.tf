resource "aws_iam_role" "ec2_ssm_role" {
  name = "go-backend-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "ec2.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "ssm_parameter_access" {
  name = "SSMParameterAccess"
  role = aws_iam_role.ec2_ssm_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:DescribeParameters"
        ],
        Resource = [
          "arn:aws:ssm:us-east-2:098917396694:parameter/app/alpha_vantage_api_key",
          "arn:aws:ssm:us-east-2:098917396694:parameter/app/cockroachdb_dsn"
        ]
      }
    ]
  })
}

resource "aws_iam_instance_profile" "ec2_instance_profile" {
  name = "go-backend-ec2-instance-profile"
  role = aws_iam_role.ec2_ssm_role.name
}