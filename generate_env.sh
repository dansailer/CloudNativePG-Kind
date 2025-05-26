#!/bin/bash

# Function to generate a random string
generate_random_string() {
    local LENGTH=$1
    LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c $LENGTH ; echo ''
}

# Generate secrets
GRAFANA_ADMIN_PASSWORD=$(generate_random_string 30)
MINIO_ROOT_PASSWORD=$(generate_random_string 30)
MINIO_USER_PASSWORD=$(generate_random_string 30)
MINIO_ACCESS_KEY=$(generate_random_string 20)
MINIO_SECRET_KEY=$(generate_random_string 40)
POSTGRES_PASSWORD=$(generate_random_string 30)
APPDBROOT_PASSWORD=$(generate_random_string 30)

# Create the .env file and write the environment variables to it
cat <<EOL > .env
# Prometheus & Grafana configuration
GRAFANA_ADMIN_PASSWORD=$GRAFANA_ADMIN_PASSWORD

# MinIO configuration
MINIO_ROOT_PASSWORD=$MINIO_ROOT_PASSWORD
MINIO_USER_PASSWORD=$MINIO_USER_PASSWORD
MINIO_ACCESS_KEY=$MINIO_ACCESS_KEY
MINIO_SECRET_KEY=$MINIO_SECRET_KEY

# PostgreSQL configuration
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
export APPDBROOT_PASSWORD=$APPDBROOT_PASSWORD
EOL

echo ".env file created successfully"
