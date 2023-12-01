FROM golang:1.21-alpine AS buildstage

COPY . /app

WORKDIR /app

RUN go install

RUN go build -o money-keeper-backend

FROM scratch

COPY --from=buildstage /app/money-keeper-backend /usr/bin

EXPOSE 8080

ENTRYPOINT [ "money-keeper-backend" ]