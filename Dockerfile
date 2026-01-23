FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY .env .env

COPY log/ log/

RUN echo "THIS IS ERORR LOGS" > ./log/error.log

RUN echo "THIS IS INFO LOGS" > ./log/info.log

COPY ./bin/proxy proxy

CMD ["./proxy"]
