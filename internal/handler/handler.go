package handler

import (
	"strconv"

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
		sizeStr := strconv.Itoa(len(valBytes))
		resp := make([]byte, 0, len(sizeStr)+2+len(valBytes))

		resp = append(resp, sizeStr...)  // ASCII digits
		resp = append(resp, '\r', '\n')  // CRLF
		resp = append(resp, valBytes...) // byte data

		return resp
	case protocol.Set:
		h.Cache.Set(cmd.Key, cmd.Value)
	}

	return []byte("2\r\nOK")
}
