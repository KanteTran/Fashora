package config

type DbPostGreSQLConfig struct {
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	DB                 string `mapstructure:"db"`
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	MaxOpenCons        string `mapstructure:"max_open_cons"`
	MaxIdleCons        string `mapstructure:"max_idle_cons"`
	ConnMaxIdleTimeSec string `mapstructure:"conn_max_idle_time_sec"`
	ConnMaxLifetimeSec string `mapstructure:"conn_max_life_time_sec"`
}

func loadDBConfig() DbPostGreSQLConfig {
	return DbPostGreSQLConfig{
		User:               GetEnv("POSTGRES_USER", "postgres"),
		Password:           GetEnv("POSTGRES_PASSWORD", "password"),
		DB:                 GetEnv("POSTGRES_DB", "postgres_db"),
		Host:               GetEnv("POSTGRES_HOST", "localhost"),
		Port:               GetEnv("POSTGRES_PORT", "5432"),
		MaxIdleCons:        GetEnv("POSTGRES_MAX_IDLE_CONS", "10"),
		MaxOpenCons:        GetEnv("POSTGRES_MAX_OPEN_CONS", "10"),
		ConnMaxLifetimeSec: GetEnv("POSTGRES_CONN_MAX_LIFE_TIME", "2000"),
		ConnMaxIdleTimeSec: GetEnv("POSTGRES_CONN_MAX_IDLE_TIME", "7"),
	}
}
