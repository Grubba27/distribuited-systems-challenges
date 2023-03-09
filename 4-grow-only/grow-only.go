package main

import (
	"context"
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	n.Handle("add", func(msg maelstrom.Message) error {
		var b map[string]any
		err := json.Unmarshal(msg.Body, &b)
		if err != nil {
			return err
		}

		r := b["delta"].(float64)
		received := int(r)
		ctx := context.TODO()
		curr, err := kv.ReadInt(ctx, "add")
		shouldCreate := false
		if err != nil {
			// does not exist
			shouldCreate = true
		}
		added := curr + received
		_ = kv.CompareAndSwap(ctx, "add", curr, added, shouldCreate)

		b["type"] = "add_ok"
		delete(b, "delta")

		return n.Reply(msg, b)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var b map[string]any
		err := json.Unmarshal(msg.Body, &b)
		if err != nil {
			return err
		}

		ctx := context.TODO()
		r, _ := kv.ReadInt(ctx, "add")
		b["value"] = r
		b["type"] = "read_ok"

		return n.Reply(msg, b)
	})

	err := n.Run()
	if err != nil {
		log.Fatal(err)
	}
}
