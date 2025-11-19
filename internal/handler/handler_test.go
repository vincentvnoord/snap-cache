package handler

import (
	"testing"

	"github.com/vincentvnoord/snap-cache/internal/cache"
	"github.com/vincentvnoord/snap-cache/internal/protocol"
)

func TestExecSetGetReturnsBytes(t *testing.T) {
	cache := cache.NewCache()
	handler := &Handler{Cache: cache}

	cmd := &protocol.Command{
		CommandType: protocol.Set,
		Key:         "key",
		Value:       []byte("value"),
	}

	handler.Exec(cmd)

	cmd.CommandType = protocol.Get
	bytes := handler.Exec(cmd)

	if string(bytes) != "value" {
		t.Fatalf("expected resulting bytes to be \"value\", but was: %s", bytes)
	}
}
