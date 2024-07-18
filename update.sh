#!/bin/bash
#
# Updates repository using docker compose to the latest release
# using github actions, see .github/workflows/deploy.yml for more
#
# Author: Kevin Gil <github.com/kevingil>

LOG_FILE="/var/log/deployment.log"

set -euo pipefail

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

handle_error() {
    log "Error occurred on line $1"
    exit 1
}

trap 'handle_error $LINENO' ERR

log "Starting deployment"

# Build and deploy
log "Building Docker images"
docker compose build || { log "Docker build failed"; exit 1; }

log "Stopping existing containers"
docker compose down || { log "Failed to stop existing containers"; exit 1; }

log "Starting new containers"
docker compose up --detach --remove-orphans || { log "Failed to start new containers"; exit 1; }

# Cleanup
log "Cleaning up unused Docker images"
docker image prune --force || log "Warning: Failed to prune Docker images"

log "Cleaning up unused Docker networks"
docker network prune --force || log "Warning: Failed to prune Docker networks"

# Check status
log "Checking container status"
docker compose ps || { log "Failed to check container status"; exit 1; }

log "Deployment completed successfully"

