package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	var b bytes.Buffer // A Buffer needs no initialization.
	b.Write([]byte("Hello \n"))

	println(b.Len())
	fmt.Fprintf(&b, "world!\n")
	println(b.Len())
	b.WriteTo(os.Stdout)
	println(b.Len())
	b.WriteTo(os.Stdout)

}

func func_1() {

	path := "./__my_reader.file"

	fd, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1024; i++ {
		fd.Write([]byte("A"))
	}

	rb := bufio.NewReaderSize(fd, 1)
	n := 0

	buffer := make([]byte, 20)
	if n, err = rb.Read(buffer); err != nil {
		panic(err)
	}
	println("Read() func return n=", n)

	if n, err = rb.Read(buffer); err != nil {
		panic(err)
	}
	println("Read() func return 2 n=", n)
}
