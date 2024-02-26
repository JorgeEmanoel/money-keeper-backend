FROM golang:1.21-alpine AS buildstage

COPY . /app

WORKDIR /app

RUN go mod tidy
RUN go build -o money-keeper-backend

FROM ubuntu:22.04

RUN mkdir /app
COPY .env.production /app/.env
COPY migrations /app/migrations
COPY --from=buildstage /app/money-keeper-backend /bin

EXPOSE 8080
