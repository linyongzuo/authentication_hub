# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . ./
COPY *.json ./
COPY *.html ./

RUN go build -o /authentication_hub

EXPOSE 8090

CMD [ "/authentication_hub" ]