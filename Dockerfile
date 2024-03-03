##########################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/app/
COPY ./src .
# Fetch dependencies.# Using go get.
RUN go get -d -v
# Build the binary.
#RUN go build -o /go/bin/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/app
############################
# STEP 2 build a small image
############################
FROM scratch
WORKDIR /
# Copy our static executable.
COPY --from=builder /go/bin/app .
COPY ./src/server.yaml .
# Run the hello binary.
ENTRYPOINT ["/app"]