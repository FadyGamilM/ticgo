FROM golang:1.21-alpine

# Set the GOPATH
ENV GOPATH /go

# Set the working directory to match your module path
WORKDIR /go/src/ticgo

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy only the auth-service source code
COPY auth-service ./auth-service

# Set working directory to auth-service inside the container  
WORKDIR /go/src/ticgo/auth-service

# Build the application
RUN go build -o auth main.go

RUN apk add --no-cache curl

EXPOSE 8080

CMD ["./auth"]