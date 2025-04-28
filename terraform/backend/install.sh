#!/bin/bash
set -x  # Enable debug output
exec > >(tee /var/log/user-data.log) 2>&1
echo "Starting install.sh at $(date)"

# Update and install necessary packages
yum update -y
yum install -y golang git
yum install -y aws-cli

# Wait for IAM profile to be available
echo "Waiting for IAM instance profile to be available..."
sleep 60

# Create app directory
mkdir -p /opt/myapp
chown ec2-user:ec2-user /opt/myapp

export AWS_REGION=us-east-2

# Verify region is set
if [ -z "$AWS_REGION" ]; then
    echo "ERROR: AWS region not configured"
    exit 1
fi

# Get API key from SSM
echo "Fetching API key from SSM..."
API_KEY=$(aws ssm get-parameter --name "/app/alpha_vantage_api_key" --with-decryption --query Parameter.Value --output text --region us-east-2)
DB_DSN=$(aws ssm get-parameter --name "/app/cockroachdb_dsn" --with-decryption --query Parameter.Value --output text --region us-east-2)

if [ -z "$API_KEY" ]; then
    echo "ERROR: Failed to retrieve API key from SSM"
    exit 1
fi
echo "API key fetched successfully"

if [ -z "$DB_DSN" ]; then
    echo "ERROR: Failed to retrieve DB DSN from SSM"
    exit 1
fi




# Create .env file in the application directory
ENV_FILE="/opt/myapp/.env"
cat <<EOF > "$ENV_FILE"
ALPHA_VANTAGE_API_KEY=$API_KEY
PORT=8080
COCKROACHDB_DSN=$DB_DSN
EOF
chown ec2-user:ec2-user "$ENV_FILE"
chmod 600 "$ENV_FILE"

# Create a systemd service that loads the .env file
cat <<EOF > /etc/systemd/system/myapp.service
[Unit]
Description=Go Backend Service
After=network.target

[Service]
ExecStart=/opt/myapp/main
WorkingDirectory=/opt/myapp
Restart=always
User=ec2-user
EnvironmentFile=/opt/myapp/.env

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable myapp.service
systemctl start myapp.service