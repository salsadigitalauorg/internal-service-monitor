package expectations_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/expectations"
	"github.com/stretchr/testify/assert"
)

func basicAuthMiddleware(t *testing.T, u string, p string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		t.Logf(auth)
		if auth == "" {
			t.Logf("No auth header")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		eu := base64.StdEncoding.EncodeToString([]byte(u+":"+p))
		if auth != "Basic "+eu {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func TestIsOK_StatusOK(t *testing.T) {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "status",
		Value: strconv.Itoa(http.StatusOK),
	}

	ok, msg := h.IsOK(c)
	if !ok {
		t.Errorf("Expected %v", http.StatusOK)
	}
	if msg != "" {
		t.Errorf("Should not return a message on success got %s", msg)
	}

}

// Ensure expected non-200 response is treated correctly.
func TestIsOK_StatusNotFound(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "pong"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "status",
		Value: strconv.Itoa(http.StatusNotFound),
	}

	ok, msg := h.IsOK(c)
	if !ok {
		t.Errorf("Expected %v", http.StatusNotFound)
	}
	if msg != "" {
		t.Errorf("Should not return a message on success got %s", msg)
	}
}

func TestIsOK_StatusNegate(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "pong"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "status",
		Value: strconv.Itoa(http.StatusOK),
		Op: "NotEqual",
	}

	ok, msg := h.IsOK(c)
	if !ok {
		t.Errorf("Expected not %v", http.StatusOK)
	}
	if msg != "" {
		t.Errorf("Should not return a message on success got %s", msg)
	}
}

func TestIsOK_StatusFailure(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "status",
		Value: strconv.Itoa(http.StatusOK),
	}

	ok, msg := h.IsOK(c)
	if ok {
		t.Errorf("Expected test to fail")
	}
	if msg != "failed" {
		t.Errorf("Expected response message 'failed' got %s", msg)
	}
}

func TestIsOK_HeaderCompare(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Header("x-ping", "pong")
		c.JSON(http.StatusNotFound, gin.H{"message": "pong"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "x-ping",
		Value: "pong",
	}

	ok, msg := h.IsOK(c)
	if !ok {
		t.Errorf("Expected %v", http.StatusNotFound)
	}
	if msg != "" {
		t.Errorf("Should not return a message on success got %s", msg)
	}
}

func TestIsOk_HeaderNotEqual(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Header("x-ping", "pong")
		c.JSON(http.StatusNotFound, gin.H{"message": "pong"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "x-ping",
		Value: "ping",
		Op: "NotEqual",
	}

	ok, msg := h.IsOK(c)
	if !ok {
		t.Errorf("Expected %v", http.StatusNotFound)
	}
	if msg != "" {
		t.Errorf("Should not return a message on success got %s", msg)
	}
}

func TestIsOk_HeaderFail(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Header("x-ping", "pong")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "header_mismatch"})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithUrl(ts.URL)

	c := cfg.MonitorExpects{
		Field: "x-ping",
		Value: "ping",
	}

	ok, msg := h.IsOK(c)
	if ok {
		t.Errorf("Expected request to fail")
	}
	if msg != "header_mismatch" {
		t.Errorf("Expected 'header_mismatch' got %s", msg)
	}
}

func TestIsOk_BasicAuth(t *testing.T) {
	r := gin.Default()
	r.Use(basicAuthMiddleware(t, "test", "test"))
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	h := expectations.Http{}
	h.WithAuth("test", "test").WithUrl(ts.URL)
	c := cfg.MonitorExpects{
		Field: "status",
		Value: strconv.Itoa(http.StatusOK),
	}

	ok, _ := h.IsOK(c)
	assert.True(t, ok)
}
