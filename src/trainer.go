package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <index_in> <db_out> <img_dir>\n", os.Args[0])
		return
	}

	ind_reader, ind_err := os.Open(os.Args[1])
	if ind_err != nil {
		log.Fatal(ind_err)
		return
	}
	defer ind_reader.Close()

	db_writer, db_err := os.Create(os.Args[2])
	if db_err != nil {
		log.Fatal(db_err)
		return
	}
	defer db_writer.Close()

	var numKanji int
	fmt.Fscanf(ind_reader, "%d", &numKanji)

	if numKanji <= 0 {
		fmt.Printf("Cannot read a number of kanji less than or equal to zero\n")
		return
	}

	fmt.Fprintf(db_writer, "%d ", numKanji)

	for i := 0; i < numKanji; i++ {
		var kanji string
		if num, err := fmt.Fscanf(ind_reader, "%s", &kanji); num == 1 && err == nil {
			img_reader, img_err := os.Open(os.Args[3] + strconv.Itoa(i) + ".png")
			if img_err != nil {
				log.Fatal(img_err)
				return
			}
			defer img_reader.Close()

			kImg, _, dec_err := image.Decode(img_reader)
			if dec_err != nil {
				log.Fatal(dec_err)
				return
			}

			fmt.Fprintf(db_writer, "%s %d ", kanji, 1)

			kVec := FeatureVector(kImg)

			for j := range kVec {
				fmt.Fprintf(db_writer, "%d ", kVec[j])
			}
		} else {
			fmt.Printf("An error occurred reading kanji from the index.\nPerhaps there weren't as many kanji as the index declared?\n")
			log.Fatal(err)
			return
		}
	}
}
