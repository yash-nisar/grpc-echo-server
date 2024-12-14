FROM golang:1.22-alpine

# Install required tools
RUN apk add --no-cache git
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server

# Add grpcurl to PATH
ENV PATH="$PATH:/root/go/bin"

# Expose the gRPC port
EXPOSE 50051

# Run the server
CMD ["./server"]