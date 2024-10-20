package main

import (
	"fmt"

	"github.com/IceCodeNew/telesend/internal/app/config"
)

func init() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("db_path: %s\n", config.TSConfig.DbPath)
	fmt.Printf("verbose: %v\n", config.TSConfig.Verbose)
}
