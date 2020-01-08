#FROM golang:alpine
#
#RUN apk add -U docker-compose
##RUN apk add -U websocket
#RUN apk add whois
#
#COPY dist/scraping /bin/
#
#RUN rm -rf dist/scraping
#
#EXPOSE 5001
#
#ENTRYPOINT [ "/bin/scraping" ]

# Start from golang base image
FROM golang:alpine as builder

# Add Maintainer info
LABEL maintainer="Daniel De La Pava Suarez <daniel·samkit@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Download all dependencies
RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/lib/pq
RUN go get github.com/go-chi/chi
RUN go get github.com/go-chi/chi/middleware
RUN go get github.com/go-chi/cors

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 5001 to the outside world
EXPOSE 5001

#Command to run the executable
CMD ["./main"]