# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Omar ElGarhy <omar.k.elgarhy@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest 

# Add /app/bin to the PATH
ENV PATH="/app/bin:${PATH}"

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./cmd/main"]