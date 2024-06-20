FROM golang:latest

WORKDIR /app

COPY . . 

RUN GOOS=linux go build -o stress-test .
