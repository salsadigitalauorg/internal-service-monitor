package expectations

import "github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"

type Expectation interface {
	IsOK(e cfg.MonitorExpects) (bool, string)
	WithUrl(u string) Expectation
}
