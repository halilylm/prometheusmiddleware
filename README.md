
## Badges

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/) [![Go Report Card](https://goreportcard.com/badge/github.com/halilylm/prometheusmiddleware)](https://goreportcard.com/report/github.com/halilylm/prometheusmiddleware)
# Prometheus Middleware For Fiber v2

Middleware for prometheus to observe metrics.

You can also register other metrics and see them in path you define.


## Metrics available by deault

`app_http_requests_total`

`app_http_request_duration_seconds`


## Installation


```bash
  go get github.com/halilylm/prometheusmiddleware@v0.0.1
```
    
## Usage/Examples

```golang
middlewarePath := "/metrics"
app := fiber.New()
registry := prometheus.NewRegistry()
middleware := prometheusmiddleware.NewPrometheusMiddleware(registry, middlewarePath)
middleware.Use(app)
```


## Feedback

If you have any feedback, please reach out to me at halilibrjim@gmail.com

