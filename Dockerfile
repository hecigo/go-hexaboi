# Building the binary of the App
FROM public.ecr.aws/q7w3q8h7/golang-1.18 AS build

WORKDIR /app

COPY go.mod go.sum ./

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download && \
    go mod verify

COPY . .

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o dora .

# Moving the binary to the 'final Image' to make it smaller
FROM public.ecr.aws/q7w3q8h7/alpine

WORKDIR /app

COPY --from=build \
    ["/app/dora", \
    "/app/favicon.ico", \
    "/app/.env", "./"]

# Exposes port 8080 because our program listens on that port
EXPOSE 8080

CMD ["./dora", "serve", "--listen", "0.0.0.0:8080"]