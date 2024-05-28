FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o chat-service

FROM busybox:latest

WORKDIR /chat-service/cmd

COPY --from=build /app/cmd/chat-service .

COPY --from=build /app/.env /chat-service

EXPOSE 8000

CMD ["./chat-service"]