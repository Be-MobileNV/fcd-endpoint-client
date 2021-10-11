package main

import (
	"math/rand"
	"time"

	ws "bitbucket.org/be-mobile/fcd-endpoint-client/client/Golang/pkg/client"
	cfg "bitbucket.org/be-mobile/fcd-endpoint-client/client/Golang/pkg/config"
	"github.com/sirupsen/logrus"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func main() {
	cfg := cfg.LoadConfig()
	logrus.Debugf("Config loaded: %s", cfg)
	wsc, err := ws.NewWebSocketClient(cfg)
	logrus.Debugf("WebSocketClient initiated: %s", wsc)
	if err != nil {
		logrus.Errorf("Could not initiate websocketclient: %s", err)
	}
	defer wsc.Close()
	for i := 0; i < 100; i++ {
		wsc.SendGPSPosition(getGPSPosition())
		time.Sleep(250 * time.Millisecond)
	}
}

func getGPSPosition() *cfg.GPSPosition {
	ymin := 46.691265
	ymax := 52.076458
	xmin := 4.565761
	xmax := 6.257655
	pos := cfg.GPSPosition{
		VehicleId:   stringWithCharset(),
		VehicleType: 1,
		EngineState: 1,
		Timestamp:   time.Now().UnixNano() / 1000000,
		Lon:         (rand.Float64() * (xmax - xmin)) + xmin,
		Lat:         (rand.Float64() * (ymax - ymin)) + ymin,
		Heading:     rand.Float32(),
		Hdop:        rand.Float32(),
		Speed:       rand.Float32(),
	}
	return &pos
}

func stringWithCharset() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
