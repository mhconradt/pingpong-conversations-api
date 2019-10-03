FROM golang:1.13-alpine3.10

# Add Maintainer Info
LABEL maintainer="Maxwell Conradt <mhconradt@protonmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/mhconradt/pingpong-conversations-api

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v -insecure ./...

RUN apk add musl-dev
RUN apk add gcc

# Install the package
RUN go install -v -a ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080
EXPOSE 10001

# Run the executable
CMD ["pingpong-conversations-api"]
