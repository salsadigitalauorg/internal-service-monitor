package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
)

type mockRouter struct {}

func (m *mockRouter) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return nil // You can mock this method as needed for your tests
}

func TestSetupRoutes(t *testing.T) {
	monitors := []cfg.MonitorConfig{
		{
			Name: "test1",
			Url: "stub|test1|test1",
			Type: "stub",
			Expects: []cfg.MonitorExpects{{Field: "test1", Value: "test1"}},
		},
		{
			Name: "test2",
			Url: "stub|test2|test2",
			Type: "stub",
			Expects: []cfg.MonitorExpects{{Field: "test2", Value: "test2"}},
		},
	}

	router := gin.Default()
	g := router.Group("/monitor")
	SetupRoutes(g, monitors)

	ts := httptest.NewServer(router)
	defer ts.Close()


	for _, monitor := range monitors {
		resp, err := http.Get(fmt.Sprintf("%s/monitor/%s", ts.URL, monitor.Name))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d; got %d", resp.StatusCode, http.StatusOK)
		}
	}
}
