FROM golang:1.22-alpine AS build
RUN apk add --no-cache git build-base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app ./

FROM alpine:latest AS runtime
RUN apk add --no-cache ca-certificates
COPY --from=build /usr/local/bin/app /usr/local/bin/app
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/usr/local/bin/app"]