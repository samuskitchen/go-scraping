# Start from golang base image
FROM golang:alpine

# Add Maintainer info
LABEL maintainer="Daniel De La Pava Suarez <daniel·samkit@gmail.com>"

# Whois is required for logic the business.
RUN apk add whois

## Install git.op
## Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create folder app
RUN mkdir /app

# Copy the source from the current directory to the working Directory inside the container
ADD . /app

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -v -o main .

# Expose port 5001 to the outside world
EXPOSE 5001

#Command to run the executable
CMD ["/app/main"]