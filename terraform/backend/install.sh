#!/bin/bash

# Update and install necessary packages
yum update -y
yum install -y golang git

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