# Linux container coming with go ROOT configured at /go.
FROM golang:1.15

# Add the source code to the container.
# We copy all source code because the explorer depends on some miner packages.
ADD . /go/src/

# CD into the source code directory.
WORKDIR /go/src/bazo-block-explorer

# Build the application.
RUN go build -o /bazo-block-explorer

# Define the start command when this container is run.
CMD ["/bazo-block-explorer", "data", ":8080", "bazo", "bazo"]