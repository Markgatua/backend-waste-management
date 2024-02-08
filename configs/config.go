package configs

import (
	"log"

	"github.com/spf13/viper"
)

func CanDebug() bool {
	return viper.GetBool("APP_DEBUG")
}

var EnvConfigs *envConfigs

// We will call this in main.go to load the env variables
func InitEnvConfigs(directory string) {
	EnvConfigs = loadEnvVariables(directory)
}

// struct to map env values
type envConfigs struct {
	AppDebug                               int     `mapstructure:"APP_DEBUG"`
	AppURL                                 string  `mapstructure:"APP_URL"`
	AfricasTalkingUsername                 string  `mapstructure:"AFRICA_TALKING_USER_NAME"`
	AfricasTalkingAPIKey                   string  `mapstructure:"AFRICA_TALKING_API_KEY"`
	DatabaseUrl                            string  `mapstructure:"DATABASE_URL"`
	JWTSecret                              string  `mapstructure:"JWT_SECRET"`
	JWTExp                                 int     `mapstructure:"JWT_EXP"`
	ArtanisVenturesPercentageCut           float64 `mapstructure:"ARTANIS_VENTURES_PERCENTAGE_CUT"`
	MailMailer                             string  `mapstructure:"MAIL_MAILER"`
	MailHost                               string  `mapstructure:"MAIL_HOST"`
	MailPort                               int     `mapstructure:"MAIL_PORT"`
	MailUsername                           string  `mapstructure:"MAIL_USERNAME"`
	MailPassword                           string  `mapstructure:"MAIL_PASSWORD"`
	MailEncryption                         string  `mapstructure:"MAIL_ENCRYPTION"`
	MailFromAddress                        string  `mapstructure:"MAIL_FROM_ADDRESS"`
	MailTo                                 string  `mapstructure:"MAIL_TO"`
	MailNoReply                            string  `mapstructure:"MAIL_NOREPLY"`
	MailFromName                           string  `mapstructure:"MAIL_FROM_NAME"`
	AccountRecoveryTokenExpirationTime     int     `mapstructure:"ACCOUNT_RECOVERY_TOKEN_EXPIRATION_TIME"`
	AcoountVerificationTokenExpirationTime int     `mapstructure:"ACCOUNT_VERIFICATION_TOKEN_EXPIRATION_TIME"`
	GeoApifyRoutePlanningApiKey            string  `mapstructure:"GEO_APIFY_ROUTE_PLANNING_API_KEY"`
}

func loadEnvVariables(directory string) (config *envConfigs) {
	viper.AddConfigPath(directory)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
