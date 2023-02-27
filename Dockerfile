FROM golang:1.19-bullseye AS build

WORKDIR /build

ENV GO111MODULE=on
ENV CGO_ENABLED=0
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o crema ./cmd/main.go


FROM alpine:3.17
RUN apk add --update --no-cache ca-certificates
WORKDIR /app

COPY --from=build /build/crema crema
COPY --from=build /build/conf.toml conf.toml

EXPOSE 8085
ENTRYPOINT ["./crema"]
