#!/bin/bash

# Start the Packulator application
cd /opt/packulator

# Pull the latest Docker image
docker pull $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/packulator:latest

# Stop and remove existing container if running
docker stop packulator-app || true
docker rm packulator-app || true

# Run the new container
docker run -d \
  --name packulator-app \
  --restart unless-stopped \
  -p 80:8080 \
  -e DB_HOST="${DB_HOST}" \
  -e DB_PORT="${DB_PORT}" \
  -e DB_USER="${DB_USER}" \
  -e DB_PASSWORD="${DB_PASSWORD}" \
  -e DB_NAME="${DB_NAME}" \
  -e APP_ENV="production" \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/packulator:latest

# Wait for application to be healthy
echo "Waiting for application to start..."
sleep 10

# Check if application is running
if curl -f http://localhost/health/check; then
    echo "Application started successfully"
    exit 0
else
    echo "Failed to start application"
    exit 1
fi