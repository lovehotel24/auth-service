/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "attempts to bring the schema for db up to date with the migrations.",
	//RunE:  migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//func migrate(cmd *cobra.Command, args []string) error {
//	cfg := database.Config{
//		User:         "postgres",
//		Password:     "postgres",
//		Host:         "localhost",
//		Name:         "users",
//		MaxIdleConns: 0,
//		MaxOpenConns: 0,
//		DisableTLS:   true,
//	}
//	db, err := database.Open(cfg)
//	if err != nil {
//		return fmt.Errorf("connect database: %w", err)
//	}
//	defer db.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	if err := schema.Migrate(ctx, db); err != nil {
//		return fmt.Errorf("migrate database: %w", err)
//	}
//
//	fmt.Println("migration complete")
//	return seed()
//}
//
//func seed() error {
//	cfg := database.Config{
//		User:         "postgres",
//		Password:     "postgres",
//		Host:         "localhost",
//		Name:         "users",
//		MaxIdleConns: 0,
//		MaxOpenConns: 0,
//		DisableTLS:   true,
//	}
//	db, err := database.Open(cfg)
//	if err != nil {
//		log.Printf("connect database: %w\n", err)
//		return fmt.Errorf("connect database: %w", err)
//	}
//	defer db.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//	defer cancel()
//
//	if err := schema.Seed(ctx, db); err != nil {
//		log.Printf("seed database:: %w\n", err)
//		return fmt.Errorf("seed database: %w", err)
//	}
//
//	log.Println("seed data complete")
//	return nil
//}
