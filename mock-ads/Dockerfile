FROM golang:1.21-alpine

WORKDIR /app

COPY main.go .

RUN go mod init mock-ads

EXPOSE 3001

CMD ["go", "run", "main.go"]