# Internal Service Monitor

This Go application monitors internal HTTP services and provides insights into their availability. It periodically connects to the configured endpoints and records the response status codes.

The intent for this application is provide HTTP endpoints that can be monitored by uptime monitors.

## Usage

``` Bash
docker run --rm -v $PWD/config.yml:/tmp/config.yml internal-service-monitor -help
-config string
  	Path to configuration file (default "cfg.yml")
-port string
  	Port to start the application on (default "8080")
-username string
  	Username for basic auth
-password string
   	Password for basic auth
```

## Installation

1. Clone the repository
``` bash
git clone https://github.com/salsadigitalauorg/internal-service-montior.git
cd internal-service-montior
```
2. Build the application
``` bash
go build -o monitor main.go
```
3. Add a configuration file
``` yaml
response_headers:
- key: x-response-header
  value: customvalue
monitors:
  - name: webhooks
    url: http://localhost:8080/ping
    type: http
    expects:
      - field: status
        value: "200"
```
4. Run the application
``` bash
$ ./monitor
```

## Configuration

The application is configured using a YAML file. You can specify the endpoints, the types of checks and what constitues a successfull check.

``` yaml
monitors:
  - name: webhooks
    url: http://localhost:8080/ping
    type: http
    expects:
      - field: status
        value: "200"
      - field: status
        value: "404"
  - name: test
    url: http://localhost:8080/test
    type: http
    expects:
      - field: status
        value: "200"
  - name: notfound
    url: http://localhost:8080/notfound
    type: http
    expects:
      - field: status
        value: "404"
  - name: httpav
    url: http://localhost:8080
    type: tcp
    expects:
      - value: "ok"

```

## License

This project is licensed under the MIT License.
