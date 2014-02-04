package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"os"
)

// The size of one side of a kanji square
const KANJI_SIDE = 64

const C_BLACK = 0
const C_WHITE = 0xFFFF

const NO_DIR = -1

var features = map[string]int {
	"WWWBBBWWW": 0, //		A		-
	"WBWWBWWBW": 1, //		B		|
	"BWWWBWWWB": 2, //		C		\
	"WWBWBWBWW": 3, //		D		/
	"WWBBBWWWW": 4, //		E		-/
	"WWWWBBBWW": 5, //		F		/-
	"BWWWBBWWW": 6, //		G		\-
	"WWWBBWWWB": 7, //		H		-\
	"WBWWBWBWW": 8, //		I		/|
	"WWBWBWWBW": 9, //		J		|/
	"WBWWBWWWB": 10, //		K		|\
	"BWWWBWWBW": 11, //		L		\|
}

/**
 * Returns whether or not a pixel in a given black-and-white image is black.
 * This function disregards the alpha value.
 */
func isBlack(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	
	return (r == g && g == b && b == C_BLACK)
}

func isWhite(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	
	return (r == g && g == b && b == C_WHITE)
}

func inBounds(img image.Image, x, y int) bool {
	bounds := img.Bounds()
	
	return	(x >= bounds.Min.X &&
			x < bounds.Max.X &&
			y >= bounds.Min.Y &&
			y < bounds.Max.Y)
}

func makeContour(img image.Image) image.Image  {
	bounds := img.Bounds()
	
	contour := image.NewRGBA(bounds)
	draw.Draw(contour, bounds, img, bounds.Min, draw.Src)
	
	// Loop 1: eliminate non-border pixels that aren't touching white
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			isBorder := false
			
			// If a pixel has a white neighbor, it's a border pixel
			for o := -1; o <= 1; o += 2 {
				if inBounds(img, x + o, y) && isWhite(img, x + o, y) {
					isBorder = true
					break
				}
				
				if inBounds(img, x, y + o) && isWhite(img, x, y + o) {
					isBorder = true
					break
				}
			}
			
			if !isBorder {
				contour.Set(x, y, color.White)
			}
		}
	}
	
	// Loop 2: eliminate corners by checking if black pixels exist horizontally and vertically
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if isBlack(contour, x, y) {
				horiz := false
				vert := false
				
				for o := -1; o <= 1; o += 2 {
					if inBounds(contour, x + o, y) && isBlack(contour, x + o, y) {
						horiz = true
					}
					
					if inBounds(contour, x, y + o) && isBlack(contour, x, y + o) {
						vert = true
					}
				}
				
				// We can eliminate this corner pixel in-place
				if horiz && vert {
					contour.Set(x, y, color.White)
				}
			}
		}
	}
	
	return contour
}

func pixDir(img image.Image, x, y int) int {
	dirarr := make([]byte, 9)
	chptr := 0
	for yo := -1; yo <= 1; yo += 1 {
		for xo := -1; xo <= 1; xo += 1 {
			if inBounds(img, x + xo, y + yo) {
				// Build up the direction string to test in the map
				if isBlack(img, x + xo, y + yo) {
					dirarr[chptr] = 'B'
				} else {
					dirarr[chptr] = 'W'
				}
				
				chptr++
			} else {
				// If we have to go out of bounds, we don't have a direction
				return NO_DIR
			}
		}
	}
	
	dirstr := string(dirarr)
	
	// If we are a feature, return our feature number
	if feature, ok := features[dirstr]; ok {
		return feature
	} else {
		return NO_DIR
	}
}

func dirMat(img image.Image) [][]int {
	bounds := img.Bounds()
	xsize := bounds.Max.X - bounds.Min.X
	ysize := bounds.Max.Y - bounds.Min.Y
	
	// Assign each pixel a direction based on its 8-neighbors
	dirs := make([][]int, ysize)
	for y := range dirs {
		dirs[y] = make([]int, xsize)
		for x := range dirs[y] {
			// Assign each pixel a direction
			// Edge pixels automatically get NO_DIR
			dirs[y][x] = pixDir(img, x, y)
		}
	}
	
	return dirs
}

func printDirMat(dirmat [][]int) {
	fmt.Print(" ")
	for i := 0; i < 64; i++ {
		fmt.Print("_")
	}
	fmt.Print(" \n\n")
	
	for y := range dirmat {
		fmt.Print("|")
		for x := range dirmat[0] {
			if dirmat[y][x] == -1 {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c", (65 + dirmat[y][x]))
			}
		}
		fmt.Print("|\n")
	}
	
	fmt.Print(" ")
	for i := 0; i < 64; i++ {
		fmt.Print("_")
	}
	fmt.Print(" ")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <img_in>", os.Args[0])
		return
	}
	
	// Create a reader and decode the data stream into an image
	reader, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	defer reader.Close()
	
	fmt.Printf("opened %s\n", os.Args[1])
	
	img, _, decerr := image.Decode(reader)
	if(decerr != nil) {
		log.Fatal(decerr)
		return
	}
	
	// Grab and print the image bounds
	bounds := img.Bounds()
	
	if bounds.Max.X != KANJI_SIDE || bounds.Max.Y != KANJI_SIDE {
		fmt.Printf("Kanji must be %dx%d\n", KANJI_SIDE, KANJI_SIDE)
		return
	}
	
	contour := makeContour(img)
	
	/** writer, werr := os.Create(os.Args[2])
	if werr != nil {
		log.Fatal(err)
		return
	}
	defer writer.Close()*/
	
	printDirMat(dirMat(contour))
}
