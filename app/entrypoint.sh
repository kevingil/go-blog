wait-for "${MYSQL_HOST}:${MYSQL_PORT}" -- "$@"

# Watch your .go files and invoke go build if the files changed.
#app-app --build="go build -o main main.go"  --command=./main
