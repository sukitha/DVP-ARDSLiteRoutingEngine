# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as build-env

# Add Maintainer Info
LABEL maintainer="Duosoftware <admin@duosoftware.com>"

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./ArdsLiteRoutingEngine/

# Create Runtime image
FROM alpine

# New Work Directory
WORKDIR /app

# Copy build and config files
COPY --from=build-env /src/main /app/

COPY --from=build-env /src/conf.json   /src/custom-environment-variables.json /app/
# Expose port 8080 to the outside world
EXPOSE 8835

# Command to run the executable
CMD ["./main"]


# Expose port 8080 to the outside world
#EXPOSE 8835

# Command to run the executable
#CMD ["./main"]
