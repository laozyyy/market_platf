FROM golang:latest

WORKDIR /app/demo
COPY . .

RUN go build big_market

EXPOSE 8080
ENTRYPOINT ["./big_market"]