##########################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"


WORKDIR $GOPATH/src/app/
COPY ./src .
# Fetch dependencies.# Using go get.
RUN go get -d -v
#RUN go mod download
# Build the binary.
#RUN go build -o /go/bin/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/app
############################
# STEP 2 build a small image
############################
FROM scratch
WORKDIR /

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group


# Copy our static executable.
COPY --from=builder /go/bin/app .
COPY ./src/server.yaml .

# Use an unprivileged user.
USER appuser:appuser
# Run the app binary.
ENTRYPOINT ["/app"]