package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"bitbucket.org/be-mobile/fcd-endpoint-client/client/Golang/pkg/config"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// // Send pings to peer with this period.
	// pingPeriod = 60 * time.Second

	// pingMessage = "fcd-endpoint-code-sample"
)

// WebSocketClient is a client that could send GPS positions over a web socket to a FCD-endpoint server.
type WebSocketClient struct {
	cfg        *config.WebSocketConfiguration
	Connection *websocket.Conn

	writeMessageLock *sync.Mutex
	Done             chan struct{}
}

// NewWebSocketClient creates a new WebSocket client.
func NewWebSocketClient(config *config.WebSocketConfiguration) (*WebSocketClient, error) {
	c := &WebSocketClient{}
	c.cfg = config
	c.writeMessageLock = &sync.Mutex{}

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

	cstDialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	header := http.Header{"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(c.cfg.Username+":"+c.cfg.Password))}}
	c.Connection, _, err = cstDialer.Dial(URL, header)
	if err != nil {
		logrus.Errorf("Dial failed: %v", err)
		return nil, err
	}

	c.Done = make(chan struct{})
	go func(wsc *WebSocketClient) {
		defer close(wsc.Done)
		for {
			_, message, err := wsc.Connection.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "close 1000 (normal)") {
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
	}(c)

	return c, err
}

// SendGPSPosition will send the GPS position to the server
func (wsc *WebSocketClient) SendGPSPosition(gpsPos *config.GPSPosition) error {
	logrus.Infof("Sending GPS position")
	// Convert to JSON string
	gpsPosJSON, err := json.Marshal(gpsPos)
	if err != nil {
		logrus.Errorf("Failed to convert the GPS position to JSON string. %v", err)
		return err
	}
	logrus.Debugf("JSON GPS postion: %s\n", gpsPosJSON)

	wsc.writeMessageLock.Lock()
	wsc.Connection.SetWriteDeadline(time.Now().Add(writeWait)) //nolint
	err = wsc.Connection.WriteMessage(websocket.TextMessage, gpsPosJSON)
	wsc.writeMessageLock.Unlock()
	if err != nil {
		logrus.Errorf("Failed to send GPS position message to server: %v", err)
		return err
	}
	return nil
}
