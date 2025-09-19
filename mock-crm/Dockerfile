FROM golang:1.21-alpine

WORKDIR /app

COPY main.go .

RUN go mod init mock-crm

EXPOSE 3002

CMD ["go", "run", "main.go"]