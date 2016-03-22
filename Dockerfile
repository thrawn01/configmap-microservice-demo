FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/thrawn01/configmap-microservice-demo

RUN go get github.com/julienschmidt/httprouter
RUN go get gopkg.in/fsnotify.v1
RUN go get gopkg.in/yaml.v2
RUN go install github.com/thrawn01/configmap-microservice-demo

# Run the command by default when the container starts.
ENTRYPOINT /go/bin/configmap-microservice-demo

# Document that the service listens on port 8080.
EXPOSE 8080
