FROM golang:1.14.0-alpine3.11 AS builder

WORKDIR /usr/local/go/src/
ADD . .

RUN go mod download
RUN go build -v -mod=readonly -o app.exe app/main.go

#lightweight docker container with binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/go/src/app.exe /
COPY --from=builder /usr/local/go/src/app.yaml /

EXPOSE 8081

CMD [ "/app.exe"]