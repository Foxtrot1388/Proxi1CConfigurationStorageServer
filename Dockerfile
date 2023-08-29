FROM golang:1.20.7-alpine3.18 AS builder

WORKDIR /usr/local/go/src/
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -mod=mod -o app.exe main.go

#lightweight docker container with binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/go/src/app.exe /
COPY --from=builder /usr/local/go/src/app.yaml /
COPY --from=builder /usr/local/go/src/*.os /
COPY --from=builder /usr/local/go/src/*.sbsl /
COPY --from=builder /usr/local/go/src/scriptcfg.json /

EXPOSE 8081

CMD [ "/app.exe"]