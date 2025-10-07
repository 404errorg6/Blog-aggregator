package config

const cfgFileName = ".gatorconfig.json"

type Config struct {
	DB_URL   string `json:"db_url"`
	UserName string `json:"current_user_name"`
}
