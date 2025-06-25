FROM golang:1.24 as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /teller cmd/teller/main.go

FROM debian:bookworm
COPY --from=build /teller /usr/bin/teller
ENTRYPOINT ["/usr/bin/teller"]