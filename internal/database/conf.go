package database

// 数据库配置
type DatabaseConfig struct {
	Driver            string `mapstructure:"driver"`
	Source            string `mapstructure:"source"`
	MaxIdleConns      int    `mapstructure:"max_idle_conns"`
	MaxOpenConns      int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime   int    `mapstructure:"conn_max_lifetime"`
	LogLevel          string `mapstructure:"log_level"`
	TablePrefix       string `mapstructure:"table_prefix"`
	SingularTable     bool   `mapstructure:"singular_table"`
	PrepareStmt       bool   `mapstructure:"prepare_stmt"`
	AllowGlobalUpdate bool   `mapstructure:"allow_global_update"`
}
