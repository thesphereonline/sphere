package p2p

import (
	"encoding/json"
	"net"
)

type MessageType string

const (
	MessageTypePing     MessageType = "ping"
	MessageTypeNewBlock MessageType = "new_block"
	MessageTypeNewTx    MessageType = "new_tx"
	MessageTypePeerList MessageType = "peer_list"
)

type Message struct {
	Type    MessageType `json:"type"`
	Payload []byte      `json:"payload"`
}

func SendMessage(conn net.Conn, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(append(data, '\n'))
	return err
}

func ReadMessage(conn net.Conn) (Message, error) {
	var msg Message
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)
	return msg, err
}
