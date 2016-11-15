package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rikuayanokozy/go-flv/flv"
)

func main() {
	var fp *os.File
	var err error

	if fp, err = os.Open(os.Args[1]); err != nil {
		panic(err)
	}

	hdr := &flv.FlvHeader{}
	if err = hdr.Parse(fp); err != nil {
		panic(err)
	}
	fmt.Println("Header:")
	fmt.Println("  Version:", hdr.Version)
	fmt.Println("  Flags:", hdr.Flags)
	fmt.Println("  Header size:", hdr.HeaderSize)
	fmt.Println("")

	i := 0
	for {
		tag := &flv.FlvTag{}
		if err = tag.Parse(fp); err == io.ErrUnexpectedEOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Printf("Tag %d:\n", i)
		fmt.Println("  PrevTagSize:", tag.PrevTagSize)
		fmt.Println("  Type:", tag.Type)
		fmt.Println("  DataSize:", tag.DataSize)
		fmt.Println("  Timestamp:", tag.Timestamp)
		fmt.Println("  TimestampEx:", tag.TimestampEx)
		fmt.Println("  StreamID:", tag.StreamID)
		if err = tag.ReadData(fp); err != nil {
			panic(err)
		}
		fmt.Println("")
		i++
	}
}
