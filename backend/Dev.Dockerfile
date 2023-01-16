FROM golang:1.19.4-alpine

WORKDIR /app

COPY . .

RUN apk add build-base

CMD [ "tail", "-f", "/dev/null" ]