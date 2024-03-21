# Shortlink Challenge

This repository hosts a simple shortlink service that provides functionality to generate short URLs for long URLs and to redirect short URLs to their corresponding URLs.

## Architecture
![image](https://github.com/1garo/shortlink/assets/44412643/592190b0-9714-40e5-840d-646a70d2aada)

## Features

- Generate short URLs for long URLs
- Redirect short URLs to their corresponding long URLs

## Getting Started

### Prerequisites

- Go installed on your machine

### Installation

1. Clone the repository:

```bash
$ git clone https://github.com/1garo/shortlink.git
$ cd shortlink
```

2. Install dependencies:

`$ go mod tidy`

3. Build the project:

`$ go run main.go`


By default, the server will start on port `8080`.

## Usage

### Generating Short URLs

To generate a short URL for a long URL, send a `POST` request to the `/shorten` endpoint with the long URL in the request body:

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com/very/long/url"}' \
  http://localhost:8080/shorten
```

The server will respond with a JSON object containing the generated short URL:

```json
{
  "short_url": "http://localhost:8080/w2fSYEJ"
}
```

### Redirecting Short URLs
To redirect a short URL to its corresponding long URL, simply visit the short URL in your browser or send a GET request to it. For example:

`curl -i http://localhost:8080/w2fSYEJ`

The server will respond with an HTTP 302 Found status code and redirect you to the long URL associated with the short URL.

## Thoughts

The app has graceful shutdown implemented, this is important when deploying our app using `k8s`. Whenever `k8s` decided to shutdown pods, it sends a `SIGTERM` signal to the application and it's important that the application is able to handle it and wait for all requests/responses to be finished and not just abruptly quits the program.

In addition to that, our app would need a `/health` route and the following configuration added to pods:
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 80
readinessProbe:
  httpGet:
    path: /health
    port: 80
```
