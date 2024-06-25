# Uso local

1. Configurar a chave do WeatherAPI no arquivo .env:
```
    KEY=xxxxxxxxxxxxxxxx
```

2. Subir o serviço:
```
    docker-compose up -d
```

3. Fazer um POST na porta 8080 em "/temp" passando o parâmetro "cep":
```
    curl -d 'cep=89222540' -i http://localhost:8080/temp
```

# Deploy no CloudRun

1. Configurar a chave do WeatherAPI no arquivo .env:
```
    KEY=xxxxxxxxxxxxxxxx
```

2. Fazer build/push/deploy no CloudRun.

# Uso do serviço rodando

Serviço rodando no Google Cloud Run:
```
    https://prod2-jph2uutbfq-rj.a.run.app
```

1. Fazer um POST em "/temp" passando o parâmetro "cep":
```
    curl -d 'cep=89222540' -i https://prod2-jph2uutbfq-rj.a.run.app/temp
```
