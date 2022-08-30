package stabilityai

import (
	"context"
	"log"
)

// Option is an option type for the client.
type Option func(*Client)

// WithContext sets the context for the client.
func WithContext(ctx context.Context) Option {
	return func(c *Client) {
		c.ctx = ctx
	}
}

// WithAPIHost sets the API host for the client.
func WithAPIHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

// WithAPIKey sets the API key for the client.
func WithAPIKey(apikey string) Option {
	return func(c *Client) {
		c.apikey = apikey
	}
}

// WithEngine sets the engine for the client.
func WithEngine(engine string) Option {
	return func(c *Client) {
		c.engine = engine
	}
}

// WithLogger sets the logger for the client.
func WithLogger(l *log.Logger) Option {
	return func(c *Client) {
		c.log = l
	}
}
