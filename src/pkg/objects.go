package pkg

type DBConfig struct {
	OrgID      string `mapstructure:"ORG_ID"`
	TenantID   string `mapstructure:"TENANT_ID"`
	UserID     string `mapstructure:"USER_ID"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBUserName string `mapstructure:"DB_USER_NAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBType     string `mapstructure:"DB_TYPE"`
	API        string `mapstructure:"API"`
	API_KEY    string `mapstructure:"API_KEY"`
}
