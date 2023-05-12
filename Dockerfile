# Build stage
FROM golang:1.20 as build

WORKDIR /build

COPY go.*  ./
RUN go mod download
COPY . .

RUN mkdir -p ./bin
RUN go generate ./...
RUN go build -o ./bin ./cmd/... 

# monitor stage
FROM gcr.io/distroless/base as monitor

WORKDIR /app

COPY --from=build /build/bin/monitor monitor
# COPY ./configs ./configs
# COPY ./LICENSE ./

EXPOSE 11029

ENTRYPOINT ["./monitor"]
