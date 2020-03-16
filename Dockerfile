
FROM golang:1.13-alpine3.11 as builder1

# Setting up environment for builder1
ENV GO111MODULE=on
WORKDIR /app/reverseproxy

# install required package(s)
RUN apk --no-cache add ca-certificates git

# Copy dependency list
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy files
COPY ./*.go ./
COPY ./proxy/*.go ./proxy/
COPY ./models/ ./models

# Compile
RUN go build -o main

# Create new stage based on alpine
FROM alpine:latest

#Copy ca certs
COPY --from=builder1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy compiled binary from builder1
WORKDIR /app
RUN mkdir /app/config/

COPY --from=builder1 /app/reverseproxy/main .

# Set Debuglevel and start the server
ENV PROXY_CONFIG "/app/config/config.toml"
ENV PROXY_DEBUG "true"
CMD [ "/app/main" ]
