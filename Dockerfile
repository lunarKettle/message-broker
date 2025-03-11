FROM golang:1.24-alpine AS build

WORKDIR /build

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o ./binary cmd/main.go

FROM gcr.io/distroless/static:latest

COPY --from=build /build/binary /app

CMD ["/app"]