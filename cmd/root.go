package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

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

	rootCmd.Flags().String("pg-user", "postgres", "user name for postgres database")
	rootCmd.Flags().String("pg-pass", "postgres", "password for postgres database")
	rootCmd.Flags().String("pg-host", "localhost", "postgres server address")
	rootCmd.Flags().String("pg-port", "5432", "postgres server port")
	rootCmd.Flags().String("pg-db", "postgres", "postgres database name")
	rootCmd.Flags().String("pg-ssl", "disable", "postgres server ssl mode on or not")
	rootCmd.Flags().String("redis-host", "127.0.0.1", "redis server host")
	rootCmd.Flags().String("redis-port", "6379", "redis server port")
	rootCmd.Flags().String("redis-user", "default", "username to access redis server")
	rootCmd.Flags().Int("redis-db", 15, "redis database name")
	rootCmd.Flags().String("redis-pass", "", "password for redis server")
	rootCmd.Flags().String("adm-ph", "0612345678", "initialize admin phone")
	rootCmd.Flags().String("adm-pass", "topSecret", "initialize admin password")
	rootCmd.Flags().String("usr-ph", "0601234567", "initialize user phone")
	rootCmd.Flags().String("usr-pass", "lowSecret", "initialize user password")
	rootCmd.Flags().String("client-id", "222222", "Oauth2 client id")
	rootCmd.Flags().String("client-secret", "22222222", "Oauth2 client secret")
	rootCmd.Flags().String("port", "8080", "auth service port")
	rootCmd.Flags().String("grpc-host", ":50051", "grpc server to connect")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("auth")
	viper.BindPFlag("pg-user", rootCmd.Flags().Lookup("pg-user"))
	viper.BindPFlag("pg-pass", rootCmd.Flags().Lookup("pg-pass"))
	viper.BindPFlag("pg-host", rootCmd.Flags().Lookup("pg-host"))
	viper.BindPFlag("pg-port", rootCmd.Flags().Lookup("pg-port"))
	viper.BindPFlag("pg-db", rootCmd.Flags().Lookup("pg-db"))
	viper.BindPFlag("pg-ssl", rootCmd.Flags().Lookup("pg-ssl"))
	viper.BindPFlag("redis-host", rootCmd.Flags().Lookup("redis-host"))
	viper.BindPFlag("redis-port", rootCmd.Flags().Lookup("redis-port"))
	viper.BindPFlag("redis-db", rootCmd.Flags().Lookup("redis-db"))
	viper.BindPFlag("redis-user", rootCmd.Flags().Lookup("redis-user"))
	viper.BindPFlag("redis-pass", rootCmd.Flags().Lookup("redis-pass"))
	viper.BindPFlag("adm-ph", rootCmd.Flags().Lookup("adm-ph"))
	viper.BindPFlag("adm-pass", rootCmd.Flags().Lookup("adm-pass"))
	viper.BindPFlag("usr-ph", rootCmd.Flags().Lookup("usr-ph"))
	viper.BindPFlag("usr-pass", rootCmd.Flags().Lookup("usr-pass"))
	viper.BindPFlag("client-id", rootCmd.Flags().Lookup("client-id"))
	viper.BindPFlag("client-secret", rootCmd.Flags().Lookup("client-secret"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("grpc-host", rootCmd.Flags().Lookup("grpc-host"))
	viper.BindEnv("gin_mode", "GIN_MODE")
	viper.AutomaticEnv()
}

func runCommand(cmd *cobra.Command, args []string) {
	dbConf := configs.NewDBConfig().
		WithHost(viper.GetString("pg-host")).
		WithPort(viper.GetString("pg-port")).
		WithUser(viper.GetString("pg-user")).
		WithPass(viper.GetString("pg-pass")).
		WithName(viper.GetString("pg-db")).
		WithSecure(viper.GetBool("pg-ssl"))

	defUser := configs.NewDefaultUser().
		WithDefaultAdminPhone(viper.GetString("adm-ph"))

	redisConf := &configs.RedisConfig{
		Addr:   fmt.Sprintf("%s:%s", viper.GetString("redis-host"), viper.GetString("redis-port")),
		DBName: viper.GetInt("redis-db"),
		Pass:   viper.GetString("redis-pass"),
		User:   viper.GetString("redis-user"),
	}

	authServerURL := fmt.Sprintf("http://localhost:%s", viper.GetString("port"))

	config := oauth2.Config{
		ClientID:     viper.GetString("client-id"),
		ClientSecret: viper.GetString("client-secret"),
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8080/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}

	clientStore := store.NewClientStore()
	clientStore.Set(viper.GetString("client-id"), &models.Client{
		ID:     viper.GetString("client-id"),
		Secret: viper.GetString("client-secret"),
		Domain: authServerURL,
	})

	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)

	db, err := configs.NewDB(dbConf)
	if err != nil {
		log.WithError(err).Error("failed to connect db")
		os.Exit(1)
	}

	err = configs.Migrate(db)
	if err != nil {
		log.WithError(err).Error("failed to migrate db schema")
	}

	userClient, err := configs.NewGrpcUserService(viper.GetString("grpc-host"))
	if err != nil {
		log.WithError(err).Error("failed to connect GRPC user service")
		os.Exit(1)
	}

	err = configs.Seed(db, defUser, userClient)
	if err != nil {
		log.WithError(err).Error("failed to seed default user")
	}

	router := gin.New()

	tokenStore := configs.NewTokenStore(redisConf)
	oauthSvr := controller.NewOauth2(db, tokenStore, clientStore)
	api := controller.NewApp(db, log, userClient, config, tokenStore, oauthSvr)
	router.Use(gin.Logger())
	routers.UserRouter(router, api)
	routers.OauthRouter(router, oauthSvr)
	if err := router.Run(fmt.Sprintf(":%s", viper.GetString("port"))); err != nil {
		log.Fatalln(err)
	}
}
