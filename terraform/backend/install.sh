#!/bin/bash

# Update and install necessary packages
yum update -y
yum install -y golang git
sudo yum install -y aws-cli

# Get API key from SSM
echo "Fetching API key from SSM..."
API_KEY=$(aws ssm get-parameter --name "/app/alpha_vantage_api_key" --with-decryption --query Parameter.Value --output text)
echo "API key fetched: $API_KEY"  # Debugging line

# Create or update .env file
ENV_FILE="/home/ec2-user/backend/.env"

# Preserve existing variables (if any) while adding/updating ALPHA_VANTAGE_API_KEY
if [ -f "$ENV_FILE" ]; then
    grep -v "ALPHA_VANTAGE_API_KEY" "$ENV_FILE" > "$ENV_FILE.tmp"
    mv "$ENV_FILE.tmp" "$ENV_FILE"
fi

echo "ALPHA_VANTAGE_API_KEY=$API_KEY" >> "$ENV_FILE"

# Set permissions
chown ec2-user:ec2-user "$ENV_FILE"

# Create a directory for the app
mkdir -p /opt/myapp
cd /opt/myapp

# Dummy file to avoid permission issues if nothing is copied yet
echo "Go backend will be deployed here" > main

# Make sure the app is executable
chmod +x main

# Create a systemd service
cat <<EOF > /etc/systemd/system/myapp.service
[Unit]
Description=Go Backend Service
After=network.target

[Service]
ExecStart=/opt/myapp/main
WorkingDirectory=/opt/myapp
Restart=always
User=ec2-user
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd, enable and start the service
systemctl daemon-reexec
systemctl daemon-reload
systemctl enable myapp.service
systemctl start myapp.service