# Stage 1: Build the binary
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.40.4

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code
COPY . .

# # Create a minimal runtime container
FROM golang:1.21-alpine

WORKDIR /app

# # Copy air executable from builder
COPY --from=builder /go/bin/air /go/bin/air

# # Copy the source code and dependencies
COPY --from=builder /app ./

ENV PORT ""

# # Expose the application's port
EXPOSE 9000

# # Use air instead of directly running the binary
# # CMD ["air", "-c", ".air.toml"]
CMD ["air"]
