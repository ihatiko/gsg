package main

import (
	"fmt"
	"gsg/cmd"
	"time"
)

func main() {
	start := time.Now()
	cmd.Run()
	end := time.Now()
	fmt.Println(end.Sub(start))
}
