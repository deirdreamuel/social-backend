FROM golang:alpine

WORKDIR /app

COPY . .
RUN go get -d -v ./...

RUN go build -o /app/bin/main cmd/main.go

EXPOSE 8080

CMD [ "/app/bin/main" ]