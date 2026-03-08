# Stage 1 - Builder
# Using official Go image to compile and run source code, naming builder to reference in multi-stage build
FROM golang:1.25-alpine AS builder

# Set Working Directory As Container, meaning creating /app and making it the current working dir in container
WORKDIR /app

# Copy dependencies from .mod & .sum into workdir into . (current working dir in container)
# We copy dependencies first so that docker layer caching - when dependencies dont change, cause image to be smaller
COPY go.mod go.sum ./

# download all dependencies from go.mod
RUN go mod download

# Copy the source code into the container
# copy all source code from current dir (.) into workdir container (.)
COPY . .

# Generate a compiled executable binary in /app called simple-product-api
RUN go build -o katalog-api

#Stage 2 - Runner
# Only using alpine, a minimal linux image without go compiler, reduce sizes
FROM alpine:latest

# Set working directory
WORKDIR /app

# Only take generated executable binary and copy it to current workdir (reduce size by removing source code, dependencies etc)
COPY --from=builder /app/katalog-api .

# Default action when running this container, run the simple-product/api binary
CMD ["./katalog-api"]




