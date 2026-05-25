package main

import (
	"context"
	"fmt"
	"os"
	"test_task_itk_academy/config"
	"test_task_itk_academy/database"
	"test_task_itk_academy/handler"
	"test_task_itk_academy/server"
)

func main() {
	cfg, err := config.Load("config.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect db: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := database.NewWalletRepository(pool)
	walletHandler := handler.NewWalletHandler(repo)
	srv := server.New(cfg.ServerPort, walletHandler)

	fmt.Printf("Server listening on :%s\n", cfg.ServerPort)
	if err := srv.Start(); err != nil {
		fmt.Println("Server error:", err)
	}
}
