package main

import (
	"conecto/factories"
	"fmt"
	"testing"
)

func TestScheduleDueSyncs(t *testing.T) {
	appConfig, _ := factories.LoadConfig[factories.AppConfig](
		"./config/conecto.json",
	)
	runner := factories.NewConecto(appConfig.ConectoConfig).Build()
	fmt.Print(runner)


}