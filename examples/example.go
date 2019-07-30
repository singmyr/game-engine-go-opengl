package main

import (
	"fmt"

	"github.com/singmyr/astrocyte"
)

func run() {
	fmt.Println("Launching astrocyte testing code")
	window, err := astrocyte.CreateWindow(640, 640, "Astrocyte testing", &astrocyte.WindowConfig{FpsLimit: 1.0})
	if err != nil {
		panic(err)
	}

	for window.IsOpen() {
		window.Update()
		window.Render()
	}
}

func main() {
	astrocyte.Run(run)
}
