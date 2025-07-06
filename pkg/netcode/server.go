package netcode

type Server struct {
	Address string
}

func NewServer(addr string) *Server {
	return &Server{Address: addr}
}
