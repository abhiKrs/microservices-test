FROM golang:1.19.4-alpine as builder

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app


# Copy the source from the current directory to the working Directory inside the container
# COPY ./app ./app
# COPY main.go .
# COPY go.mod .
COPY go.mod go.sum /

RUN go mod download

# # Download all the dependencies
# RUN go get -d -v ./...

# # Install the package
# RUN go install -v ./...

# COPY ./app ./app
# COPY main.go .
COPY . .

# RUN go build -o main .
RUN CGO_ENABLED=0 go build -o /my_main -ldflags="-w -s"
# EXPOSE 8080

# CMD [ "/app/main" ]

######## Start a new stage from scratch #######
FROM golang:1.19.4-alpine  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /my_main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./my_main"]