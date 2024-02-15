package cmd

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lovehotel24/auth-service/pkg/configs"
	"github.com/lovehotel24/auth-service/pkg/controller"
	"github.com/lovehotel24/auth-service/pkg/routers"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "auth-service",
	Short: "authentication and authorization service for love hotel24",
	Run:   runCommand,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().String("pg-user", "postgres", "user name for postgres database.")
	rootCmd.Flags().String("pg-pass", "postgres", "password for postgres database.")
	rootCmd.Flags().String("pg-host", "localhost", "postgres server address.")
	rootCmd.Flags().String("pg-port", "5432", "postgres server port.")
	rootCmd.Flags().String("pg-db", "postgres", "postgres database name.")
	rootCmd.Flags().String("pg-ssl", "disable", "postgres server ssl mode on or not.")
	rootCmd.Flags().String("redis-addr", "127.0.0.1:6379", "redis server address with port")
	rootCmd.Flags().Int("redis-db", 15, "redis database name.")
	rootCmd.Flags().String("redis-pass", "", "redis server password.")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("auth")
	viper.BindPFlag("pg-user", rootCmd.Flags().Lookup("pg-user"))
	viper.BindPFlag("pg-pass", rootCmd.Flags().Lookup("pg-pass"))
	viper.BindPFlag("pg-host", rootCmd.Flags().Lookup("pg-host"))
	viper.BindPFlag("pg-port", rootCmd.Flags().Lookup("pg-port"))
	viper.BindPFlag("pg-db", rootCmd.Flags().Lookup("pg-db"))
	viper.BindPFlag("pg-ssl", rootCmd.Flags().Lookup("pg-ssl"))
	viper.BindPFlag("redis-addr", rootCmd.Flags().Lookup("redis-addr"))
	viper.BindPFlag("redis-db", rootCmd.Flags().Lookup("redis-db"))
	viper.BindPFlag("redis-pass", rootCmd.Flags().Lookup("redis-pass"))
	viper.BindEnv("gin_mode", "GIN_MODE")
	viper.AutomaticEnv()
}

func runCommand(cmd *cobra.Command, args []string) {
	dbConf := &configs.DBConfig{
		Host:    viper.GetString("pg-host"),
		Port:    viper.GetString("pg-port"),
		User:    viper.GetString("pg-user"),
		Pass:    viper.GetString("pg-pass"),
		DBName:  viper.GetString("pg-db"),
		SSLMode: viper.GetString("pg-ssl"),
	}
	redisConf := &configs.RedisConfig{
		Addr:   viper.GetString("redis-addr"),
		DBName: viper.GetInt("redis-db"),
		Pass:   viper.GetString("redis-pass"),
	}
	router := gin.New()
	configs.Connect(dbConf)
	tokenStore := configs.NewTokenStore(redisConf)
	//sessionStore := configs.NewSessionStore()
	oauthSvr := controller.NewOauth2(configs.DB, tokenStore)
	//oauthSvr.SetPasswordAuthorizationHandler(controller.PasswordAuthorizationHandler(configs.DB))
	router.Use(gin.Logger())
	routers.UserRouter(router, oauthSvr, tokenStore)
	routers.OauthRouter(router, oauthSvr)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
