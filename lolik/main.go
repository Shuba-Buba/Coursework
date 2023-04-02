package main

import (
	"fmt"
	"lolik/conectors"
)

func main() {
	x := conectors.Connector{}
	fmt.Println(x.Connect("aboba"))
}
