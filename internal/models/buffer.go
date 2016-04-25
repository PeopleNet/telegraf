package internal_models

import (
	"github.com/influxdata/telegraf"
)

type Buffer struct {
	buf chan telegraf.Metric
	// total dropped metrics
	drops int
	// total metrics added
	total int
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		buf: make(chan telegraf.Metric, size),
	}
}

func (b *Buffer) IsEmpty() bool {
	return len(b.buf) == 0
}

func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Drops() int {
	return b.drops
}

func (b *Buffer) Total() int {
	return b.total
}

func (b *Buffer) Add(metrics ...telegraf.Metric) {
	for i, _ := range metrics {
		b.total++
		select {
		case b.buf <- metrics[i]:
		default:
			b.drops++
			<-b.buf
			b.buf <- metrics[i]
		}
	}
}

func (b *Buffer) Batch(batchSize int) []telegraf.Metric {
	n := min(len(b.buf), batchSize)
	out := make([]telegraf.Metric, n)
	for i := 0; i < n; i++ {
		out[i] = <-b.buf
	}
	return out
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}