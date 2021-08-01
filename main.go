package main

import (
	"fmt"

	"github.com/realugbun/go-config/config"
)

func main() {

	err := config.Settings.Load()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Config loaded from: " + config.Settings.ConfigLocation)
	fmt.Printf("Foo: %s, Bar: %s, Baz: %s.\n", config.Settings.Foo, config.Settings.Bar, config.Settings.Baz)
}
