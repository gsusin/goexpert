# Build

    docker build . -t stress-test

# Execução

    docker run --rm stress-test ./stress-test --url https://google.com --requests 10 --concurrency 10

    docker run --rm stress-test ./stress-test --url http://localhost:8080/course --requests 10000 --concurrency 100
