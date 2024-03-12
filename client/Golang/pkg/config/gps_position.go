package config

import (
	"fmt"
)

type GPSPosition struct {
	VehicleId   string            `json:"vehicleId"`
	VehicleType *int32            `json:"vehicleType,omitempty"`
	EngineState *int32            `json:"engineState,omitempty"`
	Timestamp   int64             `json:"timestamp"` // in unix milli
	Lon         float64           `json:"lon"`
	Lat         float64           `json:"lat"`
	Heading     *float32          `json:"heading,omitempty"` // in degrees
	Hdop        *float32          `json:"hdop,omitempty"`
	Speed       *float32          `json:"speed,omitempty"` // in km/h
	Alt         *float32          `json:"alt,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

func (g *GPSPosition) Validate() error {
	if len(g.VehicleId) > 64 {
		return fmt.Errorf("vehicleId length cannot be longer than 63, while given length is %d", len(g.VehicleId))
	}
	if g.VehicleType != nil && *g.VehicleType < 0 {
		return fmt.Errorf("invalid non-nil vehicle type:  must be positive")
	}
	if g.EngineState != nil && (*g.EngineState < -1 || *g.EngineState > 1) {
		return fmt.Errorf("invalid non-nil engine state: must be in interval [-1, 1]")
	}
	if g.Lat < -90 || g.Lat > 90 {
		return fmt.Errorf("invalid latitude: must be in interval [-90, 90]")
	}
	if g.Lon < -180 || g.Lon > 180 {
		return fmt.Errorf("invalid longitude: must be in interval [-180, 180]")
	}
	if g.Heading != nil && (*g.Heading < 0 || *g.Heading >= 360) {
		return fmt.Errorf("invalid non-nil heading: must be in interval [0, 360]")
	}
	if g.Hdop != nil && *g.Hdop < 0 {
		return fmt.Errorf("invalid non-nil hdop: must be positive")
	}
	if g.Speed != nil && *g.Speed < 0 {
		return fmt.Errorf("invalid non-nil speed: must be positive")
	}
	return nil
}
