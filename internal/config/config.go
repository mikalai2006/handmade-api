package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	vkAuthUri = "https://oauth.vk.com/authorize"
	vkRedirectUri = "http://localhost:8000/oauth/vk/me"
	VkTokenUri = "https://oauth.vk.com/access_token"
	VkUserinfoUri = "https://api.vk.com/method/users.get"
	vkScopes = []string{"account"}

	googleAuthUri = "https://accounts.google.com/o/oauth2/auth"
	googleTokenUri = "https://accounts.google.com/o/oauth2/token"
	googleUserinfoUri = "https://www.googleapis.com/oauth2/v3/userinfo"
	googleRedirectUri = "http://localhost:8000/oauth/google/me"
	googleScopes = []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}

	EnvLocal = "local"
	Prod     = "prod"

	defaultVerificationCodeLength = 10
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		Oauth OauthConfig
		// FileStorage FileStorageConfig
		// Email       EmailConfig
		// Payment     PaymentConfig
		// Limiter     LimiterConfig
		// CacheTTL    time.Duration `mapstructure:"ttl"`
		// SMTP        SMTPConfig
		// Cloudflare  CloudflareConfig
	}

	MongoConfig struct {
		Host      string
		User     string
		Password string
		Port     string
		Dbname string `mapstructure:"dbname"`
		SslMode bool `mapstructure:"sslmode"`
	}

	OauthConfig struct {
		TimeExpireCookie int

		VkAuthUri string
		VkTokenUri string
		VkUserinfoUri string
		VkClientId string
		VkClientSecret string
		VkRedirectUri string
		VkScopes []string

		GoogleAuthUri string
		GoogleTokenUri string
		GoogleUserinfoUri string
		GoogleRedirectUri string
		GoogleClientId string
		GoogleClientSecret string
		GoogleScopes []string
	}

	AuthConfig struct {
		Salt string
		SigningKey string
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`

		VerificationCodeLength int
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}
)

func Init(configsDir string) (*Config, error) {

	var cfg Config
	setDefaultConfigs(&cfg)

	// read env configs
	if err :=godotenv.Load(); err != nil {
		return nil, err
	}

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("mongodb", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}

	return viper.UnmarshalKey("oauth", &cfg.Oauth)
}

func setFromEnv(cfg *Config) {
	cfg.HTTP.Host = os.Getenv("HOST")
	cfg.HTTP.Port = os.Getenv("PORT")

	cfg.Mongo.Host = os.Getenv("MONGODB_HOST")
	cfg.Mongo.Port = os.Getenv("MONGODB_PORT")
	cfg.Mongo.User = os.Getenv("MONGODB_USER")
	cfg.Mongo.Password = os.Getenv("MONGODB_PASSWORD")

	cfg.Auth.Salt = os.Getenv("SALT")
	cfg.Auth.SigningKey = os.Getenv("SIGNING_KEY")

	cfg.Environment = os.Getenv("APP_ENV")

	cfg.Oauth.VkClientId = os.Getenv("VK_CLIENT_ID")
	cfg.Oauth.VkClientSecret = os.Getenv("VK_CLIENT_SECRET")

	cfg.Oauth.GoogleClientId = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.Oauth.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")

}

func setDefaultConfigs(cfg *Config)  {
	cfg.Oauth.VkAuthUri = vkAuthUri
	cfg.Oauth.VkRedirectUri = vkRedirectUri
	cfg.Oauth.VkScopes = vkScopes
	cfg.Oauth.VkTokenUri = VkTokenUri
	cfg.Oauth.VkUserinfoUri =VkUserinfoUri

	cfg.Oauth.GoogleAuthUri = googleAuthUri
	cfg.Oauth.GoogleRedirectUri = googleRedirectUri
	cfg.Oauth.GoogleTokenUri = googleTokenUri
	cfg.Oauth.GoogleUserinfoUri = googleUserinfoUri
	cfg.Oauth.GoogleScopes = googleScopes

	cfg.Auth.VerificationCodeLength = defaultVerificationCodeLength

}