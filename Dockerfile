# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd
RUN CGO_ENABLED=0 go build -o /go-google-photos-sync

### Deploy
FROM gcr.io/distroless/base-debian11:latest

WORKDIR /

COPY --from=build /go-google-photos-sync /go-google-photos-sync

CMD ["/go-google-photo-sync"]
