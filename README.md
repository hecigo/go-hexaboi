## Development

### Environment variables

Create `.env.local` in the project directory

```bash

# Mandatory variables
APP_NAME=hercules
APP_VERSION=v1.0.0

# (Optional) Special API root path. If it's not specified, the default is the first two characters of APP_VERSION.
APP_ROOT_PATH=v1


# ----- IN-APP CACHE -----
# View more: https://docs.gofiber.io/api/middleware/cache
# Default: false
CACHE_ENABLE=false

# Defines the prefix of cache keys. Change it to renew all of cache data.
# Full cache key is {prefix}/{path}/?{query_string}
# Optional. Default value: "hpi"
CACHE_KEY_PREFIX=hpi

# Cache default expiration duration.
# Valid time units are "ms", "s", "m", "h"
# Optional. Default value "1m"
CACHE_EXPIRATION=1m

# CacheControl enables client side caching if set to true
# Optional. Default: false
CACHE_CONTROL=false

# ETag lets caches be more efficient and save bandwidth,
# as a web server does not need to resend a full response if the content has not changed.
# Optional. Default value "true" when cache is enabled.
CACHE_ETAG=true


# ----- RESPONSE COMPRESS -----
# View more: https://docs.gofiber.io/api/middleware/compress
# Default: false
COMPRESS_ENABLE=false

# Determines the compression algorithm
# Optional. Default: 1 (LevelBestSpeed)
COMPRESS_LEVEL=1


# ----- CORS -----
# View more: https://docs.gofiber.io/api/middleware/cors
# Default: true
CORS_ENABLE=true

# Defines a list of origins that may access the resource.
# Optional. Default value "localhost"
CORS_ALLOW_ORIGINS=localhost

# Defines a list of request headers that can be used when
# making the actual request. This is in response to a preflight request.
# Optional. Default value "".
CORS_ALLOW_HEADERS=

# Defines a list methods allowed when accessing the resource.
# This is used in response to a preflight request.
# Optional. Default value "GET,POST,HEAD,PUT,PATCH"
CORS_ALLOW_METHODS=GET,POST,HEAD,PUT,PATCH

# ----- HTTP LOG -----
# View more: https://docs.gofiber.io/api/middleware/logger
# Default: false
HTTP_LOG_ENABLE=false

# Format defines the logging tags
# Optional. Default value "[${time}] | PID: ${pid} | ${latency} | ${status} ${method} ${reqHeader:X-Hpi-App-Version} ${path} ?${queryParams} ${body}"
HTTP_LOG_FORMAT=[${time}] | PID: ${pid} | ${latency} | ${status} ${method} ${reqHeader:X-Hpi-App-Version} ${path} ?${queryParams} ${body}


# ----- Pprof -----
# The handled paths all begin with /debug/pprof/.
# View more: https://docs.gofiber.io/api/middleware/pprof
# Optional. Default: false
PPROF_ENABLE=false



# ----- DB: POSTGRESQL -----
# Using GORM as default ORM library.
# View more: https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
# Inserts connection name after DB_POSTGRES if using multiple connections.
# For example, DB_POSTGRES_ABC_HOST..., "ABC" is the connection name.

# Common options
DB_POSTGRES_HOST=localhost
DB_POSTGRES_PORT=9920
DB_POSTGRES_USER=gorm
DB_POSTGRES_PWD=gorm
DB_POSTGRES_DB_NAME=gorm
DB_POSTGRES_SSL_MODE=disable
DB_POSTGRES_DSN="host=${DB_POSTGRES_HOST} user=${DB_POSTGRES_USER} password=${DB_POSTGRES_PWD} dbname=${DB_POSTGRES_DB_NAME} port=${DB_POSTGRES_PORT} sslmode=${DB_POSTGRES_SSL_MODE}"

# Sets the maximum number of connections in the idle connection pool.
# Optional. Default value 5
DB_POSTGRES_MAX_IDLE_CONNS=5

# Sets the maximum number of open connections to the database.
# Optional. Default value 20
DB_POSTGRES_MAX_OPEN_CONNS=20

# Sets the maximum amount of time a connection may be reused.
# Optional. Default value "30m"
DB_POSTGRES_CONN_MAX_LIFETIME=30m
```

### Start the application

```bash
# Via air (hot reload)
air

# Via make
make dev

# Or go run
go run main.go

# VSCode debug
Just F5 😄
```

### Use local container

```
# Clean packages
make clean-packages

# Generate go.mod & go.sum files
make requirements

# Generate docker image
make build

# Generate docker image with no cache
make build-no-cache

# Run the project in a local container
make up

# Run local container in background
make up-silent

# Stop container
make stop

# Start container
make start
```

### Use docker-compose

```
# Run docker-compose (run docker-compose with specific env file)
docker-compose --env-file .env.{*your_env_file} up
```

### Generate swagger

```
# Generate api description into swagger docs
swag init -g app.go
```

## Production

```bash
docker build -t hercules .
docker run -d -p 3000:3000 hercules
```

Go to http://localhost:3000
