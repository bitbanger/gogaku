package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 || (len(os.Args)%2) != 1 {
		fmt.Printf("Usage: %s <img_in> <img_out> ...\n")
	}

	for i := 1; i < len(os.Args); i += 2 {
		reader, err := os.Open(os.Args[i])
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer reader.Close()

		img, _, imgerr := image.Decode(reader)
		if imgerr != nil {
			log.Fatal(imgerr)
			continue
		}

		contour := MakeContour(img)

		writer, werr := os.Create(os.Args[i+1])
		if werr != nil {
			log.Fatal(werr)
			continue
		}
		defer writer.Close()

		png.Encode(writer, contour)
	}
}
