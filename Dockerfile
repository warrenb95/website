#
# Builder
#
FROM golang:1.18-alpine AS builder

# Create a workspace for the app
WORKDIR /app

# Download necessary Go modules
COPY go.mod .
RUN go mod download

# Copy over the source files
COPY src/*.go ./

# Build
RUN go build -o /main

#
# Runner
#
FROM alpine:latest AS runner

WORKDIR /

# Copy from builder the final binary
COPY --from=builder /main /main

ENTRYPOINT ["/main"]
