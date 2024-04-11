package expectations_test

import (
	"fmt"
	"net"
	"testing"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/expectations"
)

type StubTCPServer struct {
	address string
	port    int
	listener net.Listener
}

func (s *StubTCPServer) New(address string, port int) (*StubTCPServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return s, err
	}

	s.address = address
	s.port = port
	s.listener = listener

	return s, nil
}

func (s *StubTCPServer) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *StubTCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
}

func (s *StubTCPServer) Close() error {
	return s.listener.Close()
}

func TestIsOK_Success(t *testing.T) {
	s := StubTCPServer{}
	_, err := s.New("localhost", 8888)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	go s.Run()

	tcp := expectations.Tcp{}
	tcp.WithUrl("localhost:8888")

	e := cfg.MonitorExpects{Value: "ok"}

	ok, _ := tcp.IsOK(e)

	if !ok {
		t.Error("Expected to be able to connect to TCP")
	}

}

func TestIsOK_Failure(t *testing.T) {
	s := StubTCPServer{}
	_, err := s.New("localhost", 8888)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	go s.Run()

	tcp := expectations.Tcp{}
	tcp.WithUrl("localhost:10")

	e := cfg.MonitorExpects{Value: "ok"}

	ok, _ := tcp.IsOK(e)

	if ok {
		t.Error("Expected not to be able to connect to the TCP server")
	}
}
