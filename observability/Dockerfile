FROM golang:1.22-alpine
COPY . .
RUN go install ./cmd/main.go
CMD ["/go/bin/main"]

