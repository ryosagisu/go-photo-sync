# syntax=docker/dockerfile:1

## Build
FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

ENV CGO_ENABLED=0
ARG TARGETARCH
RUN GOARCH=${TARGETARCH} go build -o /go-google-photos-sync ./cmd/main.go

### Deploy
FROM alpine:3.16

WORKDIR /

COPY --from=build /go-google-photos-sync /go-google-photos-sync

CMD ["/go-google-photos-sync"]
