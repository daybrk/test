package main

import (
	"fmt"
)

func main() {
	done := Run(fmt.Sprintf(":%d", 8082))
	fmt.Println("START")
	<-done
}
