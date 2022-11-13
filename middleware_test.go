package prometheusmiddleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMiddleware(t *testing.T) {
	middlewarePath := "/metrics"
	app := fiber.New()
	registry := prometheus.NewRegistry()
	middleware := NewPrometheusMiddleware(registry, middlewarePath)
	middleware.Use(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	app.Get("/status/:code", func(c *fiber.Ctx) error {
		code, _ := strconv.Atoi(c.Params("code"))
		if code < 200 || code > 600 {
			t.Fatalf("%d is not valid http status code", code)
		}
		return c.SendStatus(code)
	})
	req := makeHttpGetRequest("/")
	res, _ := app.Test(req)
	assertStatusCode(t, http.StatusOK, res.StatusCode)
	req = makeHttpGetRequest("/status/500")
	res, _ = app.Test(req)
	assertStatusCode(t, http.StatusInternalServerError, res.StatusCode)
	req = makeHttpGetRequest("/status/404")
	res, _ = app.Test(req)
	assertStatusCode(t, http.StatusNotFound, res.StatusCode)
	req = makeHttpGetRequest(middlewarePath)
	res, _ = app.Test(req, -1)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	assertContainsString(t, string(body), `app_http_requests_total{code="200",method="GET",path="/"} 1`)
	assertContainsString(t, string(body), `app_http_requests_total{code="500",method="GET",path="/status/:code"} 1`)
	assertContainsString(t, string(body), `app_http_requests_total{code="404",method="GET",path="/status/:code"} 1`)
	assertContainsString(t, string(body), `app_http_request_duration_seconds_count{code="200"} 1`)
	assertContainsString(t, string(body), `app_http_request_duration_seconds_count{code="500"} 1`)
	assertContainsString(t, string(body), `app_http_request_duration_seconds_count{code="404"} 1`)
}

func assertStatusCode(t testing.TB, expected, got int) {
	t.Helper()
	if expected != got {
		t.Errorf("did not get correct status, expected %d, got %d", expected, got)
	}
}

func assertContainsString(t testing.TB, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("did not find the %s in %s", needle, haystack)
	}
}

func makeHttpGetRequest(path string) *http.Request {
	return httptest.NewRequest(http.MethodGet, path, nil)
}
