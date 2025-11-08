package main

import (
	"log"

	"github.com/mutition/go_start/common/config"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	
	log.Printf("%v", viper.Get("order"))
}
