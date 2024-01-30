#!/bin/bash

git pull

docker compose build

docker system prune -a

docker compose down

docker compose up -d
