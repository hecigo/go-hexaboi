## Development

### Environment variables

Create `.env.local`

```bash

# Mandatory variables
APP_NAME=gohexaboi
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



# ----- DB: POSTGRESQL, SQL Server, SQLite -----
# Using GORM as default ORM library.
# View more: https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
# Inserts connection name after DB_POSTGRES if using multiple connections.
# For example, DB_POSTGRES_ABC_HOST..., "ABC" is the connection name.
# Replace POSTGRES with SQLSERVER or SQLITE when you want to use these drivers.

# Common options
DB_POSTGRES_HOST=localhost
DB_POSTGRES_PORT=9920
DB_POSTGRES_USER=gorm
DB_POSTGRES_PWD=gorm
DB_POSTGRES_DB_NAME=gorm
DB_POSTGRES_SSL_MODE=disable

# Allows multiple connections with dot-comma (;)
DB_POSTGRES_DSN="host=${DB_POSTGRES_HOST} user=${DB_POSTGRES_USER} password=${DB_POSTGRES_PWD} dbname=${DB_POSTGRES_DB_NAME} port=${DB_POSTGRES_PORT} sslmode=${DB_POSTGRES_SSL_MODE}"
DB_POSTGRES_DSN_REPLICAS="" # Default empty string

# Sets the maximum number of connections in the idle connection pool.
# Optional. Default value 5
DB_POSTGRES_MAX_IDLE_CONNS=5

# Sets the maximum number of open connections to the database.
# Optional. Default value 20
DB_POSTGRES_MAX_OPEN_CONNS=20

# Sets the maximum amount of time a connection may be reused.
# Optional. Default value "30m"
DB_POSTGRES_CONN_MAX_LIFETIME=30m



# ----- DB: BIGQUERY -----
# As this is using the Google Cloud Go SDK, you will need to have your credentials available
# via the GOOGLE_APPLICATION_CREDENTIALS environment variable point to your credential JSON file.
# View more: https://github.com/go-gorm/bigquery
# Inserts connection name after DB_BIGQUERY if using multiple connections.
# For example, DB_BIGQUERY_ABC_HOST..., "ABC" is the connection name.

# BigQuery connection string
# Default format: bigquery://projectid/dataset
# You can also use the location format: bigquery://projectid/location/dataset
DB_BIGQUERY_DSN="bigquery://{projectid}/{location}/{dataset}"



# ----- DB: ELASTICSEARCH -----
# Elasticsearch client
# View more: https://github.com/elastic/go-elasticsearch
# Inserts connection name after ELASTICSEARCH if using multiple connections.
# For example, ELASTICSEARCH_ABC_URL..., "ABC" is the connection name.

ELASTICSEARCH_URL="localhost:9200" # Allows multiple connections with dot-comma (;)
ELASTICSEARCH_BASIC_AUTH="username:password" # Default empty string
ELASTICSEARCH_MAX_RETRIES=3 # Default 3
ELASTICSEARCH_DEBUG=false # Show debug log. Default false
ELASTICSEARCH_SEARCH_TIMEOUT=5s
ELASTICSEARCH_BATCH_INDEX_SIZE=100 # Number of docs in batch each indexing time. Default 100



# ----- DB: REDIS -----
# Redis client
# View more: https://redis.uptrace.dev/guide/universal.html

# Allows multiple connections with dot-comma (;)
# The master in the first position and then the slaves
REDIS_URL="localhost:6379"
REDIS_BASIC_AUTH="username:password" # Default empty string
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=1s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
REDIS_MASTER_NAME=""

# If you want to use cache, please add "_CACHE" to your config key
# View more: https://redis.uptrace.dev/guide/go-redis-cache.html#go-redis-cache
REDIS_CACHE_URL="localhost:6379"

# Use local in-process storage to cache the small subset of popular keys.
# These data are cleared after each deployment.
# Default cache 1000 keys for 1 minute
REDIS_CACHE_TINYFLU_SIZE=1000
REDIS_CACHE_TINYFLU_DURATION=1m

# Session
REDIS_SESSION_KEY_FORMAT="gohexaboi.session:%s/%s"
REDIS_SESSION_HASH_KEY="akjwonls"



# ----- DB: OrientDB -----
# OrientDB RESTful client
# View more: https://orientdb.org/docs/3.0.x/misc/OrientDB-REST.html
ORIENTDB_URL="localhost:2480"
ORIENTDB_DB="gohexaboi" # Database name
ORIENTDB_AUTH="username:password"
ORIENTDB_DEBUG=false # Enable debug logging. Default is false
ORIENTDB_TIMEOUT=10s # Timeout for connection. Default is 10s
ORIENTDB_MAX_RETRIES=3 # Maximum number of retry attempts. Default is 3
ORIENTDB_RETRY_WAIT_TIMEOUT=250ms # Waiting for each attempt. Default is 250ms
ORIENTDB_RETRY_MAX_WAIT_TIMEOUT=3s # Total of waiting time for attempts. Default is 3s



# ----- JWT: AWS Cognito -----
# AWS Cognito
# View more: https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html
AUTH_ENABLE=true
JWT_JWKS_URL=https://cognito-idp.{region}.amazonaws.com/{userPoolId}/.well-known/jwks.json



# ----- Elastic APM -----
ELASTIC_APM_ENABLE=true
ELASTIC_APM_SERVER_URL="localhost:8200"
ELASTIC_APM_SECRET_TOKEN="gohexaboi"
ELASTIC_APM_ENVIRONMENT=dev
ELASTIC_APM_SERVICE_NAME=$APP_NAME
ELASTIC_APM_SERVICE_VERSION=$APP_VERSION


```

### Start the application

```bash
# Via air (hot reload)
air

# Or go run
go run main.go

# VSCode debug
Just F5 😄
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
docker build -t gohexaboi .
docker run -d -p 3000:3000 gohexaboi
```

Go to http://localhost:3000
