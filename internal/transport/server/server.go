package server

import (
	"bufio"
	"errors"
	"net"
	"sync"

	"OldSchool/internal/transport/protocol"
	"OldSchool/internal/transport/router"
)

type Server interface {
	Start(port string) error
	Stop() error
}

type tcpServer struct {
	r        *router.Router
	listener net.Listener

	mu    sync.Mutex
	wg    sync.WaitGroup
	stopC chan struct{}
}

func New(r *router.Router) Server {
	return &tcpServer{
		r:     r,
		stopC: make(chan struct{}),
	}
}

func (s *tcpServer) Start(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.listener = ln
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.acceptLoop()
	}()

	return nil
}

func (s *tcpServer) Stop() error {
	s.mu.Lock()
	ln := s.listener
	s.listener = nil
	s.mu.Unlock()
	if ln != nil {
		_ = ln.Close()
	}

	close(s.stopC)
	s.wg.Wait()
	return nil
}

func (s *tcpServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stopC:
				return
			default:
			}
			return
		}

		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			_ = s.handleConn(c)
		}(conn)
	}
}

func (s *tcpServer) handleConn(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		req, err := protocol.ReadRequest(reader)
		if err != nil {
			if protocol.IsEOF(err) {
				return nil
			}
			if errors.Is(err, protocol.ErrEmptyLine) {
				_ = protocol.WriteResponse(writer, protocol.Response{
					Status:  false,
					Message: "empty request",
					Data:    nil,
				})
				continue
			}
			_ = protocol.WriteResponse(writer, protocol.Response{
				Status:  false,
				Message: "bad request",
				Data:    nil,
			})
			continue
		}

		resp := s.r.Handle(req)
		if err := protocol.WriteResponse(writer, resp); err != nil {
			return err
		}
	}
}
