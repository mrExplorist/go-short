# URL Shortener

The URL Shortener is a Go and Redis-based application that allows users to generate short and shareable links for long URLs.

## Features

- Generate short URLs from long URLs
- Redirect users to the original destination when accessing the short URL
- Efficient storage and retrieval of URLs using Redis
- Scalable and high-performance solution built with Go

## Installation

To run the URL Shortener locally, you need to have Docker and Docker Compose installed on your system. You can find installation instructions for your operating system on the [Docker website](https://docs.docker.com/get-docker/).

### Clone the Repository

```
git clone https://github.com/your-username/go-short.git
cd go-short
```

### Set up Redis

The URL Shortener uses Redis to store and retrieve URLs.

## Usage

The URL Shortener can be easily built and run using Docker Compose. Follow these steps:

1. Build and run the Docker containers:

   ```
   docker-compose up -d
   ```

   This command will build and start the URL Shortener and Redis containers in the background.

2. Access the application:

   ```
   Test the application by sending a POST request to the `http://localhost/api/v1` endpoint with a long URL in the request body:
   ```

## Configuration

The URL Shortener can be configured via environment variables in the `docker-compose.yml` file. The following variables are available:

- `DB_ADDR`= (default: db:6379): Redis server address
- `DB_PASS` (default: empty): Redis server password (if required)
- `APP_PORT` (Port on which the application will run)
- `DOMAIN` (http://localhost:APP_PORT)
  Make sure to set these variables according to your Redis configuration in the `docker-compose.yml` file before running the application.
