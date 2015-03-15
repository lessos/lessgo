package servx

type DataDriverType string

const (
	DataDriverRedis DataDriverType = "redis"
	DataDriverSSDB  DataDriverType = "ssdb"
)

type DataServiceConfig struct {
	InstanceID string `json:"instanceID"`

	// Database driver
	Driver DataDriverType `json:"driver"`

	// Exec
	Exec string `json:"exec"`

	// Database server Address. Leave blank if using unix sockets.
	Addr string `json:"host"`

	// Database server port. Leave blank if using unix sockets.
	Port uint16 `json:"port"`

	// Password for authentication.
	Pass string `json:"pass"`

	// A path of a UNIX socket file. Leave blank if using address and port.
	Socket string `json:"socket"`

	// Data Persistence Directory
	Dir string `json:"dir"`
}
