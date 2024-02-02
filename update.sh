#!/bin/bash

git pull

cd app/

docker compose build

docker system prune -a

docker compose down

docker compose up -d
