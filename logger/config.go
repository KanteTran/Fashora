package logger

// LoggerConfig for console mode, json format, log level
type LoggerConfig struct {
	EnableConsole     bool   `json:"console" mapstructure:"console"`
	ConsoleJSONFormat bool   `json:"format_json" mapstructure:"format_json"`
	ConsoleLevel      string `json:"level" mapstructure:"level"`
}

func LoggerDefaultConfig() LoggerConfig {
	return LoggerConfig{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      "info",
	}
}
