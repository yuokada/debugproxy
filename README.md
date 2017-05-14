# debugproxy
Proxy server for request debug

## Install

``` sh
% go get github.com/yuokada/debugproxy
```

## Usage

``` sh
% debugproxy -h
Usage of debugproxy:
  -drequest
        set true if you debug Request (default true)
  -dresponse
        set true if you debug Response (default true)
  -dst string
        proxy destination (default "http://localhost:8080")
  -port int
        listen port (default 8081)
```

Client side

``` bash
% curl -i http://localhost:8081/
```

Server side

``` bash
% debugproxy -dst http://localhost:8000
------- Debug start --------
GET / HTTP/1.1
Host: localhost:8081
Accept: */*
User-Agent: curl/7.51.0
X-Forwarded-For: ::1

------- Debug end --------
------- Debug start --------
HTTP/1.0 200 OK
Connection: close
Content-Length: 19
Content-Type: text/html
Date: Sun, 14 May 2017 06:05:38 GMT
Last-Modified: Sun, 14 May 2017 02:32:45 GMT
Server: SimpleHTTP/0.6 Python/2.7.13

<h1>It Works!</h1>
------- Debug end --------
```
