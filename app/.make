.PHONY: pull-rebuild

pull-rebuild:
    docker-compose down
    git pull
    docker-compose build
    docker-compose up -d
