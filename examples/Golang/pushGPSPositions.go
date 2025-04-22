package main

import (
	"context"
	"encoding/json"
	"time"

	ws "github.com/Be-MobileNV/fcd-endpoint-client/client/Golang/pkg/client"
	"github.com/Be-MobileNV/fcd-endpoint-client/client/Golang/pkg/config"
	"github.com/sirupsen/logrus"
)

// Load config, create WebSocketClient and send 100 random GPS positions
// to the specified endpoint.
func main() {
	cfg := &config.WebSocketConfiguration{
		Address: "127.0.0.1",
		Port:    "8080",
		TLS:     false,
	}
	wsc, err := ws.NewWebSocketClient(context.Background(), cfg)
	if err != nil {
		logrus.Errorf("Could not initiate websocketclient: %v", err)
		return
	}
	defer wsc.Close()
	for range 100 {
		pos := getGPSPosition()
		p, _ := json.Marshal(pos)
		logrus.Infof("Sending %d bytes", len(p))
		err = wsc.SendGPSPosition(context.Background(), pos)
		if err != nil {
			logrus.Errorf("failed to send GPS position: %v", err)
		}
		time.Sleep(250 * time.Millisecond)
	}
}
