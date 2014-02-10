package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <img_in> <db_in>\n", os.Args[0])
		return
	}

	img_reader, img_err := os.Open(os.Args[1])
	fmt.Println(os.Args[1])
	if img_err != nil {
		log.Fatal(img_err)
		return
	}
	defer img_reader.Close()

	dbr, db_err := os.Open(os.Args[2])
	if db_err != nil {
		log.Fatal(db_err)
		return
	}
	defer dbr.Close()

	vecdb := make(map[string][][]int)

	numkanji := 0
	fmt.Fscanf(dbr, "%d", &numkanji)

	for i := 0; i < numkanji; i++ {
		var kanji string

		fmt.Fscanf(dbr, "%s", &kanji)

		var numvec int
		fmt.Fscanf(dbr, "%d", &numvec)

		vecvec := make([][]int, numvec)

		for j := 0; j < numvec; j++ {
			vecvec[j] = make([]int, 196)

			for k := 0; k < 196; k++ {
				fmt.Fscanf(dbr, "%d", &vecvec[j][k])
			}
		}

		vecdb[kanji] = vecvec
	}

	kImg, _, err := image.Decode(img_reader)
	if err != nil {
		log.Fatal(err)
		return
	}

	kVec := FeatureVector(kImg)

	fmt.Printf("input character looks like %s\n", KanjiClass(kVec, vecdb))
}
