package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Be-MobileNV/fcd-endpoint-client/client/Golang/pkg/config"
	"github.com/coder/websocket"
	"github.com/sirupsen/logrus"
)

// Time allowed to write a message to the peer.
const writeWait = 10 * time.Second

// WebSocketClient is a client that could send GPS positions over a web socket to a FCD-endpoint server.
type WebSocketClient struct {
	cfg        *config.WebSocketConfiguration
	Connection *websocket.Conn
}

// NewWebSocketClient creates a new WebSocket client.
func NewWebSocketClient(ctx context.Context, config *config.WebSocketConfiguration) (*WebSocketClient, error) {
	c := WebSocketClient{cfg: config}

	URL := fmt.Sprintf("wss://%s:%s/v1/ws", c.cfg.Address, c.cfg.Port)
	if !c.cfg.TLS {
		URL = fmt.Sprintf("ws://%s:%s/v1/ws", c.cfg.Address, c.cfg.Port)
	}

	_, err := url.Parse(URL)
	if err != nil {
		logrus.Errorf("URL parsing failed: %v", err)
		return nil, err
	}

	logrus.Infof("Connecting to %s", URL)

	header := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(c.cfg.Username+":"+c.cfg.Password))}}
	opts := &websocket.DialOptions{
		HTTPHeader:      header,
		CompressionMode: websocket.CompressionContextTakeover,
	}
	c.Connection, _, err = websocket.Dial(ctx, URL, opts)
	if err != nil {
		logrus.Errorf("Dial failed: %v", err)
		return nil, err
	}

	go func(wsc *WebSocketClient) {
		for {
			_, message, err := wsc.Connection.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					logrus.Infof("Received message from server: '%v'", err)
				} else {
					logrus.Errorf("Failed to read the message from the server: %v", err)
				}
				return
			}
			// a normal message indicates an error from the endpoint (e.g. an input parsing error) if it starts with "ERR: "
			if strings.HasPrefix(string(message), "ERR: ") {
				logrus.Errorf("Received error message from server: %s", message)
			} else {
				logrus.Infof("Received message from server: %s", message)
			}
		}
	}(&c)

	return &c, err
}

// SendGPSPosition will send the GPS position to the server
func (wsc *WebSocketClient) SendGPSPosition(ctx context.Context, gpsPos *config.GPSPosition) error {
	if err := gpsPos.Validate(); err != nil {
		return fmt.Errorf("validation of gpsPosition did not succeed: %w", err)
	}
	logrus.Infof("Sending GPS position")
	// Convert to JSON string
	gpsPosJSON, err := json.Marshal(gpsPos)
	if err != nil {
		logrus.Errorf("Failed to convert the GPS position to JSON string. %v", err)
		return err
	}
	logrus.Debugf("JSON GPS postion: %s\n", gpsPosJSON)

	writeCtx, writeCancel := context.WithTimeout(ctx, writeWait)
	defer writeCancel()
	err = wsc.Connection.Write(writeCtx, websocket.MessageText, gpsPosJSON)
	if err != nil {
		logrus.Errorf("Failed to send GPS position message to server: %v", err)
		return err
	}
	return nil
}

// Close will send a close message to the server
func (wsc *WebSocketClient) Close() {
	logrus.Infof("Closing the websocket by sending a close message")
	err := wsc.Connection.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		logrus.Errorf("Failed to send close message to server: %v", err)
	}
}
