package http

import (
	"context"
	"github.com/dmalykh/refurbedsender/sender"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpGate_send(t *testing.T) {

	type config struct {
		timeout time.Duration
	}
	type server struct {
		responseCode int
		sleep        time.Duration
	}

	tests := []struct {
		name    string
		config  config
		server  server
		wantErr bool
	}{
		{
			name: `Send http OK`,
			config: config{
				timeout: 1 * time.Second,
			},
			server: server{
				responseCode: 200,
			},
			wantErr: false,
		},
		{
			name: `Send http with empty timeout — context deadline`,
			server: server{
				responseCode: 200,
			},
			wantErr: true,
		},
		{
			name: `Send http Not found`,
			config: config{
				timeout: 1 * time.Second,
			},
			server: server{
				responseCode: 404,
			},
			wantErr: true,
		},
		{
			name: `Send http timeout`,
			config: config{
				timeout: 1 * time.Second,
			},
			server: server{
				responseCode: 200,
				sleep:        5 * time.Second,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		func() {
			var srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.server.sleep)
				w.WriteHeader(tt.server.responseCode)
				_, _ = w.Write([]byte(`something`)) //nolint:errcheck
			}))
			defer func() {
				time.Sleep(tt.server.sleep)
				srv.Close()
			}()

			t.Run(tt.name, func(t *testing.T) {
				var h = &Gate{
					url:     srv.URL,
					timeout: tt.config.timeout,
				}

				if err := h.Send(context.TODO(), sender.NewMessage([]byte{})); (err != nil) != tt.wantErr {
					t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}()
	}
}
