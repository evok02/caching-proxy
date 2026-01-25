FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY ./log/ ./log/

COPY ./configs/ ./configs/

COPY ./bin/proxy proxy

CMD ["./proxy"]
