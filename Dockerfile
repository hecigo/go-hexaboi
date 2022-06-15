# Building the binary of the App
FROM golang:1.18 AS build

WORKDIR /go/src/hoangphuc.tech

# Copy all the Code and stuff to compile everything
COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download \
    && go mod verify

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o gohexaboi .

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest

WORKDIR /app

COPY --from=build \
    ["/go/src/hoangphuc.tech/go-hexaboi", \
    "/go/src/hoangphuc.tech/favicon.ico", \
    "/go/src/hoangphuc.tech/.env", "./"]

# Exposes port 8080 because our program listens on that port
EXPOSE 8080

CMD ["./gohexaboi","localhost:8080"]