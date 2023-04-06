# Use an official Golang runtime as a parent image
FROM arm64v8/golang:1.20-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install required dependencies
RUN go mod download

# Build the program
RUN go build -o main .

# Expose port 8080 for the program to listen on
EXPOSE 8080

# Run the program when the container starts
CMD ["./main"]
