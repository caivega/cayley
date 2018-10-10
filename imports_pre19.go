// +build !go1.9

package cayley

import (
	"github.com/caivega/cayley/graph"
	"github.com/caivega/cayley/graph/path"
)

type Iterator graph.Iterator
type QuadStore graph.QuadStore
type QuadWriter graph.QuadWriter

type Path path.Path

type Handle struct {
	graph.QuadStore
	graph.QuadWriter
}

func (h *Handle) Close() error {
	err := h.QuadWriter.Close()
	h.QuadStore.Close()
	return err
}
