package peer

import (
	"goraft/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Peer struct {
	ID     int64
	Addr   string
	Conn   *grpc.ClientConn
	Client proto.NodeClient
}

func (p *Peer) EnsureConnected() error {
	if p.Conn != nil && p.Client != nil {
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

func (p *Peer) Close() error {
	if p.Conn != nil {
		return p.Conn.Close()
	}

	return nil
}

func NewPeer(id int64, addr string) Peer {
	return Peer{
		ID:   id,
		Addr: addr,
	}
}
