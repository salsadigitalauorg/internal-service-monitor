package expectations

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
)


type Tcp struct {
	Url string
	Port int
}

func (t *Tcp) WithUrl(u string) Expectation {
	if strings.HasPrefix(u, "http://") {
		u = strings.TrimPrefix(u, "http://")
	} else if strings.HasPrefix(u, "https://") {
		u = strings.TrimPrefix(u, "https://")
	}

	parts := strings.Split(u, ":")

	if len(parts) != 2 {
		log.Printf("Parts is not long enough...")
		return t
	}

	t.Url = parts[0]
	t.Port, _ = strconv.Atoi(parts[1])

	return t
}

func (t *Tcp) IsOK(e cfg.MonitorExpects) (bool, string) {
	dial, err := net.Dial("tcp", fmt.Sprintf("%s:%v", t.Url, t.Port))
	if err != nil {
		log.Printf("Err: %v", err.Error())
		return false, err.Error()
	}
	defer dial.Close()
	return true, ""
}
