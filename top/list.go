package top

import (
	"fmt"
	"io"

	"github.com/migotom/heavykeeper"
)

// Item is top List single item.
type Item struct {
	Name  string
	Count uint64
}

// List is list of top items.
type List []Item

// WriteResults prints List into nice table.
func (l List) WriteResults(w io.Writer) error {
	for _, item := range l {
		if _, err := fmt.Fprintln(w, item.Name, "\t", item.Count); err != nil {
			return err
		}
	}

	// Hack for tabwriter
	if v, ok := w.(interface{ Flush() error }); ok {
		return v.Flush()
	}

	return nil
}

// GetListFromHK return List from heavykeeper's []minheap.Node.
// So we can forget about heavykeeper and work with it`s results.
func GetListFromHK(k uint32, hk *heavykeeper.TopK) List {
	t := make(List, 0, k)
	for _, v := range hk.List() {
		// The top has fewer values than requested
		if v.Item == "" {
			break
		}

		t = append(t, Item{
			Name:  v.Item,
			Count: v.Count,
		})
	}

	return t
}
