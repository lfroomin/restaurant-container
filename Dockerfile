# Use an official Golang runtime as a parent image
FROM arm64v8/golang:1.20-alpine as builder

# Set the working directory to /app
WORKDIR /app

# Install required dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the current directory contents into the container at /app
COPY . /app

# Build the program
RUN go build -o main .

FROM arm64v8/alpine
COPY --from=builder /app/main /main
COPY --from=builder /app/app.env /app.env

# Expose port 8080 for the program to listen on
EXPOSE 8080

# Set the working directory to /
WORKDIR /

# Run the program when the container starts
CMD ["./main"]
