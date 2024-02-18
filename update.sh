#!/bin/bash

git pull

cd app/

docker compose build

docker compose down

docker compose up -d

docker system prune -a
