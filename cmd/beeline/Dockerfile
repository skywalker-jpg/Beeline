FROM ubuntu:latest

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY ./bin/beeline/beeline /app/beeline
COPY ./config/config.yaml /app

CMD ["/app/beeline"]