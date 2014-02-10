package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <db_out> <img_dir>\n", os.Args[0])
		return
	}

	db_writer, db_err := os.Create(os.Args[1])
	if db_err != nil {
		log.Fatal(db_err)
		return
	}
	defer db_writer.Close()

	files, fil_err := ioutil.ReadDir(os.Args[2])
	if fil_err != nil {
		log.Fatal(fil_err)
		return
	}

	numKanji := len(files)

	fmt.Fprintf(db_writer, "%d ", numKanji)

	for i := range files {
		file := files[i].Name()

		img_reader, img_err := os.Open(path.Join(os.Args[2], file))
		if img_err != nil {
			log.Fatal(img_err)
			return
		}

		kImg, _, dec_err := image.Decode(img_reader)
		if dec_err != nil {
			log.Fatal(dec_err)
			return
		}

		kanji := strings.Split(file, ".")[0]

		// The 1 is a placeholder till I get a more varied dataset
		fmt.Fprintf(db_writer, "%s %d ", kanji, 1)

		kVec := FeatureVector(kImg)

		for j := range kVec {
			fmt.Fprintf(db_writer, "%d ", kVec[j])
		}

		img_reader.Close()
	}
}
