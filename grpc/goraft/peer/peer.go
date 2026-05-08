package peer

import (
	"goraft/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Peer represents a peer in the Raft cluster, containing its ID, address, gRPC connection, and client interface.
type Peer struct {
	ID     int64
	Addr   string
	Conn   *grpc.ClientConn
	Client proto.NodeClient
}

// Open establishes a gRPC connection to the peer's address and initializes the client interface.
func (p *Peer) Open() error {
	if p.Conn != nil {
		return nil
	}

	conn, err := grpc.NewClient(
		p.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return err
	}

	p.Conn = conn
	p.Client = proto.NewNodeClient(conn)

	return nil
}

// Close terminates the gRPC connection to the peer if it exists.
func (p *Peer) Close() error {
	if p.Conn != nil {
		return p.Conn.Close()
	}

	return nil
}

// NewPeer creates a new Peer with the given ID and address.
func NewPeer(id int64, addr string) Peer {
	return Peer{
		ID:   id,
		Addr: addr,
	}
}
