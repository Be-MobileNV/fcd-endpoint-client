/*
 * fcd-endpoint-client
 *
 * FCD-endpoint-client
 *
 * API version: 1.0.0
 * Contact: api-support@be-mobile.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// GpsPosition GPS Position
type GpsPosition struct {
	VehicleId   string            `json:"vehicleId"`
	VehicleType int32             `json:"vehicleType"`
	EngineState int32             `json:"engineState"`
	Timestamp   int32             `json:"timestamp"`
	Lon         float32           `json:"lon"`
	Lat         float32           `json:"lat"`
	Heading     float32           `json:"heading"`
	Hdop        float32           `json:"hdop"`
	Speed       float32           `json:"speed"`
	Metadata    map[string]string `json:"metadata"`
}
