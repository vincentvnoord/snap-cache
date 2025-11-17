package handler

import (
	"github.com/vincentvnoord/internal/cache"
	"github.com/vincentvnoord/internal/protocol"
)

type Handler struct {
	Cache *cache.Cache
}

func (h *Handler) Exec(cmd *protocol.Command) []byte {
	switch cmd.CommandType {
	case protocol.Ping:
		return []byte("PONG")
	case protocol.Get:
		val := h.Cache.Get(cmd.Key)
		if val == nil {
			return []byte("")
		}

		return val.Value
	case protocol.Set:
		h.Cache.Set(cmd.Key, cmd.Value)
	}

	return []byte{}
}
