package config

type WebSocketConfiguration struct {
	Address  string `required:"true" flagUsage:"The address of the server."`
	Port     string `default:"443" flagUsage:"The port of the server."`
	Username string `required:"true" flagUsage:"The username of the basic authorization."`
	Password string `required:"true" flagUsage:"The password of the basic authorization."`
	TLS      bool   `default:"false" flagUsage:"Use secure communication"`
}

type GPSPosition struct {
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
