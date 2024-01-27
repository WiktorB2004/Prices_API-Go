# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Prepare go projectO
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the project files
COPY . . 
# Build the Go app
RUN go build -o /app/bin/main .

# Expose port 8080 to the outside world
EXPOSE 3001

# Command to run the executable
CMD ["/app/bin/main"]

