# Build stage
FROM golang:1.20 AS build

WORKDIR /build

COPY go.*  ./
RUN go mod download
COPY . .

RUN mkdir -p ./bin
RUN go generate ./... \
 && CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o ./bin ./cmd/...

# monitor stage
FROM scratch AS monitor

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/bin/monitor /

ENTRYPOINT ["/monitor"]
