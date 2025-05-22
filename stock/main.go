package main

import (
	_ "common/config"
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	serviceName := viper.Get("stock.service-name")
	fmt.Println(serviceName)
}
