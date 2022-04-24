package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dmalykh/refurbedsender/gate"
	"github.com/dmalykh/refurbedsender/sender"
	"net/http"
	"net/url"
	"time"
)

// Gate implements the gate.Gate interfaces for sending requests via http standard library
type Gate struct {
	url     string
	timeout time.Duration
}

func (h *Gate) Send(ctx context.Context, message sender.Message) error {
	return h.send(ctx, message)
}

// Send request with default HTTP client and return error
func (h *Gate) send(ctx context.Context, message sender.Message) error {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url, bytes.NewReader(message.GetText()))
	if err != nil {
		return err
	}
	var client = http.DefaultClient
	client.Timeout = h.timeout

	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(`response code should by %d, got %d`, http.StatusOK, resp.StatusCode)
	}

	return err
}

type Config struct {
	URL     string
	Timeout time.Duration
}

func NewHTTPGate(c *Config) (gate.Gate, error) {
	val, err := url.Parse(c.URL)
	if err != nil {
		return nil, err
	}
	var g = &Gate{
		url:     val.String(),
		timeout: c.Timeout,
	}
	if g.timeout == 0 {
		g.timeout = 30 * time.Second
	}

	return g, nil
}
