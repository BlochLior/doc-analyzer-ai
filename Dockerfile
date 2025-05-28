# Stage 1: Build the Go application
# Use the official Go image with a specific, recent stable version
FROM golang:1.24.3-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
# This allows Go to download dependencies before copying the source code,
# which helps with Docker layer caching. If go.mod/go.sum don't change,
# this layer won't be rebuilt.
COPY go.mod go.sum ./

# Download Go module dependencies
# This command fetches all required modules and stores them in the module cache.
RUN go mod download

# Copy the entire Go source code to the working directory
# The `.` at the end means "copy everything from the current directory on the host to /app in the container".
COPY . .

# Build the Go application
# -o specifies the output file name ('server' in this case)
# ./cmd/server is the path to your main package
# -ldflags="-s -w" reduces the binary size by stripping debug information
# CGO_ENABLED=0 is important for static binaries, making the final image smaller and more portable.
# GOOS=linux ensures the binary is built for a Linux environment, compatible with Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: Create the final lean image
# Use a minimal base image for the final application (alpine is very small)
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /app

# Copy the compiled executable from the 'builder' stage to the final image
# '/app/server' refers to the path where the executable was built in the previous stage
COPY --from=builder /app/server .

# Expose port 8080 so that the container can receive traffic on this port
# This is documentation for Docker, it doesn't actually publish the port.
EXPOSE 8080

# Command to run when the container starts
# This executes your compiled Go application
CMD ["./server"]