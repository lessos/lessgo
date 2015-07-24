package base

type Config struct {
	// Database driver
	Driver string `json:"driver"`

	// Database server hostname or IP. Leave blank if using unix sockets.
	Host string `json:"host"`

	// Database server port. Leave blank if using unix sockets.
	Port string `json:"port"`

	// Username for authentication.
	User string `json:"user"`

	// Password for authentication.
	Pass string `json:"pass"`

	// A path of a UNIX socket file. Leave blank if using host and port.
	Socket string `json:"socket,omitempty"`

	// Name of the database.
	Dbname string `json:"dbname,omitempty"`

	// Database engine.
	Engine string `json:"engine,omitempty"`

	// Database charset.
	Charset string `json:"charset,omitempty"`

	// Data Source Name
	Dsn string `json:"dsn,omitempty"`
}
