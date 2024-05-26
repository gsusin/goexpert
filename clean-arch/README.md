# Passos para executar a aplicação

Iniciar os servidores do MySQL e do RabbitMQ:

    docker-compose up -d

Criar o banco de dados:

    make migrate

Executar a aplicação:

    cd cmd/ordersystem

    go run main.go wire_gen.go

Os serviços responderão nas portas:
- web server: 8000
- gRPC server: 50051
- GraphQL server: 8080



