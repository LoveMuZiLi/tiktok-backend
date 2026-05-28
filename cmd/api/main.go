package main

import (
	"log"

	"github.com/LoveMuZiLi/tiktok-backend/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
