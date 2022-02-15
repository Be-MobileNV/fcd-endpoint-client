package config

import (
	"flag"
	"os"
	"unsafe"

	"github.com/koding/multiconfig"
	"github.com/sirupsen/logrus"
)

type WebSocketConfiguration struct {
	LogLevel string `default:"info"`
	Address  string `required:"true" flagUsage:"The address of the server."`
	Port     string `default:"443" flagUsage:"The port of the server."`
	Username string `required:"true" flagUsage:"The username of the basic authorization."`
	Password string `required:"true" flagUsage:"The password of the basic authorization."`
	TLS      bool   `default:"false" flagUsage:"Use secure communication"`
}

type GPSPosition struct {
	VehicleId   string            `json:"vehicleId"`
	VehicleType int32             `json:"vehicleType,omitempty"`
	EngineState int32             `json:"engineState,omitempty"`
	Timestamp   int64             `json:"timestamp"`
	Lon         float64           `json:"lon"`
	Lat         float64           `json:"lat"`
	Heading     float32           `json:"heading,omitempty"`
	Hdop        float32           `json:"hdop,omitempty"`
	Speed       float32           `json:"speed,omitempty"`
	Alt         float32           `json:"alt,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

func (g *GPSPosition) Validate() bool {
	if unsafe.Sizeof(g.VehicleId) > 64 {
		return false
	}
	if g.VehicleType < 0 || g.VehicleType > 19 {
		return false
	}
	if g.EngineState < -1 || g.EngineState > 1 {
		return false
	}
	if g.Lat < -90 || g.Lat > 90 {
		return false
	}
	if g.Lon < -180 || g.Lon > 180 {
		return false
	}
	if g.Heading < 0 || g.Heading >= 360 {
		return false
	}
	if g.Hdop < 0 {
		return false
	}
	if g.Speed < 0 {
		return false
	}
	return true
}

// LoadConfig reads the configuration
func LoadConfig() *WebSocketConfiguration {
	m := multiconfig.New()
	cfg := &WebSocketConfiguration{}
	err := m.Load(cfg)
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		logrus.Fatalf("Failed to load config: %+v", err)
	} else {
		logrus.Infof("Loaded cfg %+v", cfg)
	}
	if err := m.Validate(cfg); err != nil {
		logrus.Fatalf("Invalid config: %+v", err)
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatalf("Invalid log level %s : %+v", cfg.LogLevel, err)
	} else {
		logrus.SetLevel(lvl)
	}

	return cfg
}
