package main

import (
	"flag"
	"fmt"
)

func main() {
	webPort := flag.Int("p", 8082, "порт веб сервера")
	done := Run(fmt.Sprintf(":%d", *webPort))
	fmt.Println("START")
	<-done
}
