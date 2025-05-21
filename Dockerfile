FROM golang:1.23.0-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o node ./cmd/node

FROM scratch

WORKDIR /

COPY --from=builder /build/node .

COPY data /data

ENTRYPOINT ["./node"]
