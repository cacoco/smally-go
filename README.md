# smal.ly
[![Run Tests](https://github.com/cacoco/smally-go/actions/workflows/ci.yml/badge.svg)](https://github.com/cacoco/smally-go/actions/workflows/ci.yml)

A simple url shortener written in go+redis using [echo](https://echo.labstack.com/guide/). 

## Usage

## Testing

```
$ go test ./...
```

### Running

To run locally using foreman:

```
$ PORT=8080 REDIS_HOST=localhost REDIS_PORT=6379 foreman start web
16:05:48 web.1  | started with pid 60842
16:05:48 web.1  | 
16:05:48 web.1  |    ____    __
16:05:48 web.1  |   / __/___/ /  ___
16:05:48 web.1  |  / _// __/ _ \/ _ \
16:05:48 web.1  | /___/\__/_//_/\___/ v4.7.2
16:05:48 web.1  | High performance, minimalist Go web framework
16:05:48 web.1  | https://echo.labstack.com
16:05:48 web.1  | ____________________________________O/_______
16:05:48 web.1  |                                     O\
16:05:48 web.1  | ⇨ http server started on [::]:8080
```

### Using

To create a new shortened url: post a JSON body to the `/url` endpoint in the form of `{"url":"TO_SHORTEN_URL"}`

```
$ curl -i -H "Content-Type: application/json" -X POST -d '{"url":"http://www.nytimes.com/2012/05/06/travel/36-hours-in-barcelona-spain.html"}' http://127.0.0.1:8080/url
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Content-Length: 44

{"smally_url":"http://127.0.0.1:8080/9h5k4"}
```

Then in a browser paste the shortened URL to be redirected to the original URL or use [curl](http://curl.haxx.se/docs/manual.html):

```
$ curl -i http://127.0.0.1:8080/9h5k4
HTTP/1.1 301 Moved Permanently
Location: http://www.nytimes.com/2012/05/06/travel/36-hours-in-barcelona-spain.html
Content-Length: 0
```