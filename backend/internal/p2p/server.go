package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/thesphereonline/sphere/internal/core/block"
	"github.com/thesphereonline/sphere/internal/core/chain"
	"github.com/thesphereonline/sphere/internal/db"
)

type P2PServer struct {
	Port     string
	Peers    []*Peer
	Chain    *chain.Blockchain
	Database *db.Postgres
}

func NewP2PServer(port string) *P2PServer {
	return &P2PServer{
		Port:  port,
		Peers: []*Peer{},
	}
}

func (s *P2PServer) Start() error {
	ln, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}
	fmt.Println("🔗 P2P server listening on port", s.Port)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn)
		}
	}()

	return nil
}

func (s *P2PServer) ConnectToPeer(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	peer := &Peer{Conn: conn}
	s.Peers = append(s.Peers, peer)
	fmt.Println("🤝 Connected to peer:", address)
	return nil
}

func (s *P2PServer) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("⚠️ Error reading from peer:", err)
			return
		}

		raw = strings.TrimSpace(raw)
		var incoming block.Block
		if err := json.Unmarshal([]byte(raw), &incoming); err != nil {
			fmt.Println("❌ Invalid block from peer:", err)
			continue
		}

		// Check if we already have this block
		if s.Chain.Contains(incoming.Hash) {
			continue
		}

		fmt.Println("📥 New block received from peer:", incoming.Hash)

		// Save to DB
		if err := s.Database.SaveBlock(&incoming); err != nil {
			fmt.Println("❌ Failed to persist incoming block:", err)
			continue
		}

		s.Chain.AddBlock(&incoming)
	}
}
