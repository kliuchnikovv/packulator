#!/bin/bash

# Stop the Packulator application
echo "Stopping Packulator application..."

# Stop the Docker container
docker stop packulator-app || true

# Remove the Docker container
docker rm packulator-app || true

echo "Application stopped successfully"