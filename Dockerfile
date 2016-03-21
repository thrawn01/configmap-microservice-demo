# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/thrawn01/configmap-microservice-demo

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/julienschmidt/httprouter
RUN go get gopkg.in/fsnotify.v1
RUN go get gopkg.in/yaml.v2
RUN go install github.com/thrawn01/configmap-microservice-demo

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/configmap-microservice-demo

# Document that the service listens on port 8080.
EXPOSE 8080
