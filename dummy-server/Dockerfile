# syntax=docker/dockerfile:1
FROM golang:1.17
WORKDIR /go/src/github.com/marcubus/vrising-server/
COPY dummy.go ./
COPY go.mod ./
COPY go.sum ./
RUN GOOS=linux go build -o dummy

FROM debian:bookworm-slim
RUN apt update
RUN apt install -y ca-certificates
EXPOSE 9876/udp
WORKDIR /root/
COPY --from=0 /go/src/github.com/marcubus/vrising-server/dummy ./
CMD ["./dummy"]
