package expectations

import (
	"strings"

	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
)

type Expectation interface {
	IsOK(e cfg.MonitorExpects) (bool, string)
	WithUrl(u string) Expectation
	WithAuth(u string, p string) Expectation
}

type Stub struct {
	Url string
	Username string
	Password string
	ShouldReturn string
	ShouldMsg string
}

func (t *Stub) WithUrl(u string) Expectation {
	s := strings.Split(u, "|")
	t.Url = s[0]
	t.ShouldReturn = s[1]
	t.ShouldMsg = s[2]
	return t
}

func (t *Stub) WithAuth(u string, p string) Expectation {
	t.Username = u
	t.Password = p
	return t
}

func (t *Stub) IsOK(e cfg.MonitorExpects) (bool, string) {
	if e.Value == t.ShouldReturn {
		return true, t.ShouldMsg
	} else {
		return false, "nomatch"
	}
}
