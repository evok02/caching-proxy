FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY .env .env

COPY ./bin/proxy proxy

CMD ["./proxy"]
