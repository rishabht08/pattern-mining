FROM golang:1.22.2 as builder
WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# RUN go build -o main .

# Use a minimal Debian-based image to run the Go application
# FROM gcr.io/distroless/base-debian10

# Set the Current Working Directory inside the container
# WORKDIR /

# Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/main .

# Command to run the executable
# CMD ["./main"]
CMD ["go", "run", "src/main.go"]