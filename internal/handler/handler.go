package handler

import (
	"github.com/vincentvnoord/snap-cache/internal/cache"
	"github.com/vincentvnoord/snap-cache/internal/protocol"
)

type Handler struct {
	Cache *cache.Cache
}

// HAS TO write \r\n to the response so the client knows what to read
func (h *Handler) Exec(cmd *protocol.Command) []byte {
	switch cmd.CommandType {
	case protocol.Ping:
		return []byte("PONG\r\n")
	case protocol.Get:
		entry := h.Cache.Get(cmd.Key)
		if entry == nil {
			return []byte("\r\n")
		}

		valBytes := entry.Value
		resp := make([]byte, 0, len(valBytes)+2)
		resp = append(resp, valBytes...)
		resp = append(resp, '\r', '\n')

		return resp
	case protocol.Set:
		h.Cache.Set(cmd.Key, cmd.Value)
	}

	return []byte("OK\r\n")
}
