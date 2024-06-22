# Uso

Aplicação rodando no Google Cloud Run:

    https://prod2-jph2uutbfq-rj.a.run.app

Fazer um POST em "/temp" passando o parâmetro "cep":

    curl -d 'cep=89222540' https://prod2-jph2uutbfq-rj.a.run.app/temp

# Deploy no CloudRun

1. Configure a chave do WeatherAPI no arquivo .env:

    KEY=xxxxxxxxxxxxxxxx

2. Fazer build/push/deploy no CloudRun.
