package database

// 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Source          string `mapstructure:"source"`
	DatabaseName    string `mapstructure:"database_name"`
	MaxPoolSize     uint64 `mapstructure:"max_pool_size"`
	MinPoolSize     uint64 `mapstructure:"min_pool_size"`
	MaxIdleTime     int    `mapstructure:"max_idle_time"`
	LogLevel        string `mapstructure:"log_level"`
}
