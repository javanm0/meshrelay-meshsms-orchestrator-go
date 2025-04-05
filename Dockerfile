# Use the official Go image as the base image for building the application
FROM golang:1.24-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go module file
COPY go.mod ./

# Create an empty go.sum file if it doesn't exist
RUN touch go.sum

# Download dependencies (if any)
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a smaller base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /app/main .

# Expose the port the app runs on
EXPOSE 3040

# Command to run the application
CMD ["./main"]