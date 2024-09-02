package main

import (
	"fmt"
	"jps/internal/config"
)

func main() {
	cfg := config.MustReadConfig()
	fmt.Print(cfg)
}
