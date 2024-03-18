package main

import (
	"math/rand"
	"time"

	cfg "bitbucket.org/be-mobile/fcd-endpoint-client/client/Golang/pkg/config"
)

// characters to be used in the VehicleID
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func ptr[T any](x T) *T { return &x }

// generate random GPS Position to push to endpoint
func getGPSPosition() *cfg.GPSPosition {
	ymin := 46.691265
	ymax := 52.076458
	xmin := 4.565761
	xmax := 6.257655
	pos := cfg.GPSPosition{
		VehicleId:   stringWithCharset(),
		VehicleType: ptr[int32](1),
		EngineState: ptr[int32](1),
		Timestamp:   time.Now().UnixMilli(),
		Lon:         (rand.Float64() * (xmax - xmin)) + xmin,
		Lat:         (rand.Float64() * (ymax - ymin)) + ymin,
		Heading:     ptr(rand.Float32()),
		Hdop:        ptr(rand.Float32()),
		Speed:       ptr(rand.Float32() * 120),
	}
	return &pos
}

// generate random VehicleID
func stringWithCharset() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
