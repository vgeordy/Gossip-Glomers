package main

import (
	"encoding/json"
	"log"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastResponse struct {
    Type string `json:"type"`
}

type TopologyResponse struct {
    Type string `json:"type"`
}

type ReadResponse struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}


func main() {
	n := maelstrom.NewNode()

	var list []int

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages := int(body["message"].(float64))
		list = append(list, messages)

		return n.Reply(msg, BroadcastResponse {Type: "broadcast_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resp := ReadResponse{
			Type:     "read_ok",
			Messages: list,
		}

		return n.Reply(msg, resp)

	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		return n.Reply(msg, TopologyResponse {Type: "topology_ok"})
	})

		if err := n.Run(); err != nil {
			log.Fatal(err)
		}
}
