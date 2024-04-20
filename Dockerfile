FROM golang:1.22 as build
WORKDIR /usr/src/app
RUN apt-get update && apt-get install libx11-dev -y
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /teller cmd/teller/main.go

FROM debian:bookworm
COPY --from=build /teller /usr/bin/teller
ENTRYPOINT ["/usr/bin/teller"]