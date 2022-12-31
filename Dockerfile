# syntax=docker/dockerfile:1

## Build
FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

ENV CGO_ENABLED=0
ARG TARGETARCH
RUN GOARCH=${TARGETARCH} go build -o /go-photo-sync ./cmd/main.go

### Deploy
FROM alpine:3.16

WORKDIR /

COPY --from=build /go-photo-sync /go-photo-sync

ENV PATH="$PATH:/go-photo-sync"
ENV CONFIG_PATH=""
ENV IMAGE_PATH=""
ENTRYPOINT ["/go-photo-sync"]
CMD ["-command=SyncImage"]
