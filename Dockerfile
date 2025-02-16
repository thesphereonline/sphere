FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /blockchain-node ./cmd/node

# Create data directory
RUN mkdir /data

EXPOSE 8080

ENTRYPOINT ["/blockchain-node"] 