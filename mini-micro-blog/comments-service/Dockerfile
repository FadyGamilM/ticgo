# Stage 1: Build the binary
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Stage 2: Create a minimal runtime container
FROM alpine:latest

# Set the working directory in the runtime container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

ENV EVENT_BUS_SVC_PORT ""
ENV PORT ""

# Expose the application's port
EXPOSE 9001

# Set the default command
CMD ["./app"]
