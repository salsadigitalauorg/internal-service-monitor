package expectations

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
)

type Http struct {
	Url string
	Status int
	Header http.Header
	Username string
	Password string
	Body JsonBody
	ResponseTime time.Duration
	Error string
}

type JsonBody struct {
	Message string `json:"message"`
}

func (h *Http) WithUrl(u string) Expectation {
	start := time.Now()
	h.Url = u

	req, _ := http.NewRequest("GET", h.Url, nil)
	if h.Username != "" && h.Password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(h.Username + ":" + h.Password))
		req.Header.Add("Authorization", "Basic " + auth)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return h
	}

	h.Status = res.StatusCode
	h.Header = res.Header

	defer res.Body.Close()

	h.ResponseTime = time.Since(start)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return h
	}

	var b JsonBody
	json.Unmarshal(body, &b)
	h.Body = b

	return h
}

func (h *Http) WithAuth(u string, p string) Expectation {
	h.Username = u
	h.Password = p
	return h
}

func (h *Http) IsOK(c cfg.MonitorExpects) (bool, string) {
	if c.Field == "status" {
		if h.IsExpectedStatus(c.Value, c.Op) {
			return true, ""
		} else {
			return false, h.Body.Message
		}
	} else {
		if h.HasHeaderWithValue(c.Field, c.Value, c.Op) {
			return true, ""
		} else {
			return false, h.Body.Message
		}
	}
}

func (h *Http) IsExpectedStatus(e string, op string) bool {
	expected, _ := strconv.Atoi(e)

	switch op {
		case "Equal", "Eq", "":
			return h.Status == expected
		case "NotEqual", "Ne":
			return h.Status != expected
	}

	return false
}

func (h *Http) HasHeaderWithValue(key string, value string, op string) bool {
	header := h.Header.Get(key)

	switch op {
		case "Equal", "Eq", "":
			return header == value
		case "NotEqual", "Ne":
			return header != value
	}

	return false
}
