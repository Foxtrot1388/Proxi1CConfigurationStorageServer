FROM golang:1.14.0-alpine3.11 AS builder

WORKDIR /usr/local/go/src/
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -v -mod=readonly -o app.exe main.go

#lightweight docker container with binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/go/src/app.exe /
COPY --from=builder /usr/local/go/src/app.yaml /
COPY --from=builder /usr/local/go/src/*.os /

EXPOSE 8081

CMD [ "/app.exe"]