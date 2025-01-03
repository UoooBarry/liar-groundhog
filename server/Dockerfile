# Use the slim version of Go for both build and runtime
FROM golang:1.23.4-alpine AS base

# Set the working directory
WORKDIR /app

# Copy the entire project folder into the container
COPY . .

# Expose the application's port
EXPOSE 8080

RUN go mod download

# Production
FROM base AS prod
RUN go build -o ./dist/liar-groundhog
CMD ["./dist/liar-groundhog"]

# Development
FROM base AS dev
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]
