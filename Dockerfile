#
# Builder
#
FROM golang:1.18-alpine AS builder

# Create a workspace for the app
WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy over everything
COPY . ./

# Build
RUN go build -o /main

#
# Runner
#
FROM alpine:latest AS runner

# Copy from builder the final binary
COPY --from=builder /main /

EXPOSE 8080

ENTRYPOINT ["/main"]
