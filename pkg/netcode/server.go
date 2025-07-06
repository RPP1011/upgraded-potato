package netcode

import (
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	Address  string
	upgrader websocket.Upgrader

	listener net.Listener
	srv      *http.Server
	mu       sync.Mutex
	conns    map[*websocket.Conn]struct{}
}

func NewServer(addr string) *Server {
	s := &Server{
		Address:  addr,
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		conns:    make(map[*websocket.Conn]struct{}),
	}
	return s
}

// Start begins listening for websocket connections on the server's address.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}
	s.listener = ln
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)
	srv := &http.Server{Handler: mux}
	s.srv = srv
	go srv.Serve(ln)
	return nil
}

// Stop shuts down the server and closes all active connections.
func (s *Server) Stop() error {
	if s.srv != nil {
		_ = s.srv.Close()
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.mu.Lock()
	for c := range s.conns {
		c.Close()
	}
	s.conns = make(map[*websocket.Conn]struct{})
	s.mu.Unlock()
	return nil
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	if len(s.conns) >= 3 {
		s.mu.Unlock()
		http.Error(w, "connection limit reached", http.StatusServiceUnavailable)
		return
	}
	s.mu.Unlock()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.mu.Lock()
	s.conns[conn] = struct{}{}
	s.mu.Unlock()

	go s.readLoop(conn)
}

func (s *Server) readLoop(c *websocket.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.conns, c)
		s.mu.Unlock()
		c.Close()
	}()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		if mt != websocket.BinaryMessage {
			continue
		}
		var m Message
		if err := proto.Unmarshal(msg, &m); err != nil {
			break
		}
		resp, err := proto.Marshal(&m)
		if err != nil {
			break
		}
		if err := c.WriteMessage(websocket.BinaryMessage, resp); err != nil {
			break
		}
	}
}
