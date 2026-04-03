package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"mom/cmd"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("[aviso] .env nao carregado:", err)
	}


	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
