FROM golang:latest

RUN mkdir -p /app
WORKDIR /app
ADD . /app

CMD {"src .env"}

RUN go mod vendor
RUN go build .

ENTRYPOINT [ "/app/stubber" ]
