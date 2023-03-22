#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

RUN addgroup -g 10014 choreo && \
    adduser  --disabled-password --uid 10014 --ingroup choreo choreouser

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app

USER 10014

ENTRYPOINT /app
LABEL Name=gofilesystembenchmark Version=0.0.1
EXPOSE 8080
