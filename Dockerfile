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

# Download all dependencies
RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/lib/pq
RUN go get github.com/go-chi/chi
RUN go get github.com/go-chi/chi/middleware
RUN go get github.com/go-chi/cors

# Build the Go app
RUN go build -v -o main .

# Expose port 5001 to the outside world
EXPOSE 5001

#Command to run the executable
CMD ["/app/main"]