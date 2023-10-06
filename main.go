package main

import (
	"fmt"
	"io"
)

func main() {
	c := N
	pg := pcapgo.NewWriter(out)
	f := gopacket.NewSerializeBuffer().Bytes()
	fmt.Println(f)
}
