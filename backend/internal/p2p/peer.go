package p2p

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/thesphereonline/sphere/internal/core/block"
)

type Peer struct {
	Conn net.Conn
}

func (p *Peer) SendBlock(b *block.Block) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(p.Conn, "%s\n", data)
	return err
}
