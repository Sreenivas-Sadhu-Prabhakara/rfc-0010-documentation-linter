package lab

import (
	"context"
	"fmt"
)

const MaxPayload = 1006

type Frame struct {
	Link        int    `json:"link"`
	Destination int    `json:"destination"`
	Sequence    int    `json:"sequence"`
	Payload     []byte `json:"payload"`
}

// Pipeline is our modern framing model, not the historical binary header.
// The bounded channel transfers ownership of fresh payload copies to the handler.
func Pipeline(ctx context.Context, data []byte, link, dest, capacity int, send func(Frame) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if link < 0 || link > 255 || dest < 0 || dest > 31 || capacity < 1 || capacity > 1024 || send == nil {
		return fmt.Errorf("invalid pipeline configuration")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	queue := make(chan Frame, capacity)
	done := make(chan error, 1)
	go func() {
		defer close(queue)
		seq := 0
		for pos := 0; pos < len(data); pos += MaxPayload {
			end := pos + MaxPayload
			if end > len(data) {
				end = len(data)
			}
			f := Frame{link, dest, seq, append([]byte(nil), data[pos:end]...)}
			seq++
			select {
			case <-ctx.Done():
				done <- ctx.Err()
				return
			case queue <- f:
			}
		}
		done <- nil
	}()
	var result error
	for f := range queue {
		if result != nil {
			continue
		}
		if err := ctx.Err(); err != nil {
			result = err
			cancel()
			continue
		}
		if err := send(f); err != nil {
			result = err
			cancel()
		}
	}
	producerErr := <-done
	if result != nil {
		return result
	}
	if producerErr != nil {
		return producerErr
	}
	return ctx.Err()
}
