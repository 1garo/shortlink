# Shortlink Challenge

This repository hosts a simple shortlink service that provides functionality to generate short URLs for long URLs and to redirect short URLs to their corresponding URLs.

## Architecture
![image](https://github.com/1garo/shortlink/assets/44412643/51658204-2b36-4700-b3a9-c3c405e08a2e)

## Features

- Generate short URLs for long URLs
- Redirect short URLs to their corresponding long URLs

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Make](https://www.gnu.org/software/make/#download)
- [Nginx](https://www.nginx.com/resources/wiki/start/topics/tutorials/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. Clone the repository:

```bash
$ git clone git@github.com:1garo/shortlink.git
$ cd shortlink
```

2. Up all needed containers:

```bash
$ make up
```

- Remove all containers:

  ```bash
  $ make down
  ```


By default on dev, the server will start on port `3000`, you can use `make run`.

On  docker compose, nginx runs on port `9999` and forward traffic to both server running.

Tests can be run with `make test` or `make testv` for verbose output.

Checkout [Makefile](./Makefile) to see all possible commands.

## Usage

### Generating Short URLs

To generate a short URL for a long URL, send a `POST` request to the `/shorten` endpoint with the long URL in the request body:

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url"}' \
  http://localhost:9999/shorten
```

The server will respond with a JSON object containing the generated short URL:

```json
{
  "short_url": "ZQy5lnI"
}
```

### Redirecting Short URLs
To redirect a short URL to its corresponding long URL, simply visit the short URL in your browser or send a GET request to it. For example:

`curl -i http://localhost:9999/ZQy5lnI`

or just use the URL on your favorite browser (I really encourage you guys to go this URL :D).

The server will respond with an `HTTP 302 Found` and redirect you to the long URL associated with the short URL.

## Load Balancer
I choose `nginx` because I feel that is the most simple and effective (and the one that I have some experience).

Load balancer configuration is always tricky, I created this [one](nginx.conf) with rate limit in mind but trying to not be so strict and remove the service capability of processing requests.

## Database
The choice was to go with `mongodb` because it fits better our requirements of horizontal auto-scaling, for example, with the possibility of using `shards`(way easier to setup than SQL). 

`shortUrl` field is using a `text-index` because it's our `key` in all operations.

## Improvements
1. The app has graceful shutdown implemented, a must feature when using `kubernetes`. Whenever `kubernetes` decided to shutdown pods, it sends a `SIGTERM` signal to the application and it's important that the application is able to handle it and wait for all requests/responses to be finished and not just abruptly quits the program.

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
2. Auth: Could be a fancy implementation like OAuth 2.0 or a simpler one like JWT.
3. I always like to let infrastructure do most of the hardwork related to scaling, performance and optimizations. But after that we can go to the code and start doing some improvements, like introducing goroutines.
4. A cache layer you be great(e.g Redis), would improve our performance and our services would be free to process new incoming requests.
