package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	n := maelstrom.NewNode()
	var msgs []int
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var b map[string]any
		err := json.Unmarshal(msg.Body, &b)
		if err != nil {
			return err
		}

		i := b["message"].(int)
		msgs = append(msgs, i)

		b["type"] = "broadcast_ok"
		delete(b, "message")

		return n.Reply(msg, b)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var b map[string]any
		err := json.Unmarshal(msg.Body, &b)
		if err != nil {
			return err
		}
		b["messages"] = msgs
		b["type"] = "read_ok"
		return n.Reply(msg, b)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {

		var b map[string]any
		err := json.Unmarshal(msg.Body, &b)
		if err != nil {
			return err
		}
		delete(b, "topology")
		b["type"] = "topology_ok"

		return n.Reply(msg, b)
	})

	err := n.Run()
	if err != nil {
		log.Fatal(err)
	}
}
