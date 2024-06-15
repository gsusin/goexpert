# Configuração

Configure no arquivo .env:

- chave do WheatherAPI

    KEY=xxxxxxxxxxxxxxxx

# Execução

1. Faça o build e execute a aplicação:

    docker-compose up --build --detach zipkin-collector zipkin-client
    
2. Faça um POST em http://localhost:8080/temp, passando um CEP válido no parâmetro "cep"

3. Verifique nos logs o traceId:

    docker-compose logs --tail=1 zipkin-client

4. Abra o Zipkin em 

    http://localhost:9411/zipkin/traces/{traceId}

5. Verifique os tempos de resposta do serviceA e serviceB.
