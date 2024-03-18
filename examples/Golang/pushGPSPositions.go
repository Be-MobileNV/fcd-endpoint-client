package main

import (
	"time"

	ws "bitbucket.org/be-mobile/fcd-endpoint-client/client/Golang/pkg/client"
	cfg "bitbucket.org/be-mobile/fcd-endpoint-client/client/Golang/pkg/config"
	"github.com/sirupsen/logrus"
)

// load config, create websocketclient and send 100 random gpspositions
// to the specified endpoint.
func main() {
	cfg := cfg.LoadConfig()
	logrus.Debugf("Config loaded: %v", cfg)
	wsc, err := ws.NewWebSocketClient(cfg)
	logrus.Debugf("WebSocketClient initiated: %v", wsc)
	if err != nil {
		logrus.Errorf("Could not initiate websocketclient: %v", err)
	}
	defer wsc.Close()
	for i := 0; i < 100; i++ {
		wsc.SendGPSPosition(getGPSPosition())
		time.Sleep(250 * time.Millisecond)
	}
}
