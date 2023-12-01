FROM golang:1.21.4-bullseye

COPY . /app

WORKDIR /app

RUN go install

RUN go build -o money-keeper-backend
RUN cp money-keeper-backend /usr/bin

EXPOSE 8080

ENTRYPOINT [ "money-keeper-backend" ]