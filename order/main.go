package main

import (
	_ "common/config" // import for side effect to load configuration
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	serviceName := viper.GetString("order.service-name")
	fmt.Println(serviceName)
}
