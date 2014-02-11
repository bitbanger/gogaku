package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"math"
)

const C_WHITE = 0xFFFF

const NO_DIR = -1

// Each feature key in this map describes a 3x3 pixel area from left to right and then top down.
var features = map[string]int{
	"WWWBBBWWW": 0,  //		A		-
	"WBWWBWWBW": 1,  //		B		|
	"BWWWBWWWB": 2,  //		C		\
	"WWBWBWBWW": 3,  //		D		/
	"WWBBBWWWW": 4,  //		E		-/
	"WWWWBBBWW": 5,  //		F		/-
	"BWWWBBWWW": 6,  //		G		\-
	"WWWBBWWWB": 7,  //		H		-\
	"WBWWBWBWW": 8,  //		I		/|
	"WWBWBWWBW": 9,  //		J		|/
	"WBWWBWWWB": 10, //		K		|\
	"BWWWBWWBW": 11, //		L		\|
}

func isBlack(img image.Image, x, y int) bool {
	return !isWhite(img, x, y)
}

func isWhite(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()

	return (r == g && g == b && b == C_WHITE)
}

func inBounds(img image.Image, x, y int) bool {
	bounds := img.Bounds()

	return (x >= bounds.Min.X &&
		x < bounds.Max.X &&
		y >= bounds.Min.Y &&
		y < bounds.Max.Y)
}

// MakeContour takes a normalized kanji and removes all non-border pixels, 
// producing a contour line image.
func MakeContour(img image.Image) image.Image {
	bounds := img.Bounds()

	contour := image.NewRGBA(bounds)
	draw.Draw(contour, bounds, img, bounds.Min, draw.Src)

	// Loop 1: eliminate non-border pixels that aren't touching white
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			isBorder := false

			// If a pixel has a white neighbor, it's a border pixel
			for o := -1; o <= 1; o += 2 {
				if inBounds(img, x+o, y) && isWhite(img, x+o, y) {
					isBorder = true
					break
				}

				if inBounds(img, x, y+o) && isWhite(img, x, y+o) {
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
					if inBounds(contour, x+o, y) && isBlack(contour, x+o, y) {
						horiz = true
					}

					if inBounds(contour, x, y+o) && isBlack(contour, x, y+o) {
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

// pixDir determines the line direction of an individual pixel.
func pixDir(img image.Image, x, y int) int {
	dirarr := make([]byte, 9)
	chptr := 0
	for yo := -1; yo <= 1; yo += 1 {
		for xo := -1; xo <= 1; xo += 1 {
			if inBounds(img, x+xo, y+yo) {
				// Build up the direction string to test in the map
				if isBlack(img, x+xo, y+yo) {
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

 	// Some versions of Go need this?
	return -1
}

// dirMat produces a matrix of integers describing the line direction of each pixel.
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

// printDirMat prints a matrix of characters that describe the line direction of pixels.
func printDirMat(dirmat [][]int) {
	fmt.Print(" ")
	for i := 0; i < 64; i++ {
		fmt.Print("_")
	}
	fmt.Print(" ")

	for y := range dirmat {
		fmt.Print("|")
		for x := range dirmat[0] {
			if dirmat[y][x] == -1 {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c", (65 + dirmat[y][x]))
			}
		}
		fmt.Print("|")
	}

	fmt.Print(" ")
	for i := 0; i < 64; i++ {
		fmt.Print("_")
	}
	fmt.Print(" ")
}

// These range functions do not work properly on their own.
// However, they can classify ranges exclusively when applied in outward order (see below).

func inARange(x, y int) bool {
	return (x >= 6 && x <= 9) && (y >= 6 && y <= 9)
}

func inBRange(x, y int) bool {
	return (x >= 4 && x <= 11) && (y >= 4 && y <= 11)
}

func inCRange(x, y int) bool {
	return (x >= 2 && x <= 13) && (y >= 2 && y <= 13)
}

// FeatureVector extracts a 196-dimensional vector that describes a 64x64 kanji.
func FeatureVector(img image.Image) []int {
	// Convert our image into contour form and then into a directional matrix
	contour := MakeContour(img)
	dirmat := dirMat(contour)

	// Initialize our feature counters
	// One for each exclusive sub-window
	// These are used for each sub-window, but we allocate them here and reset them in the loop
	aFeats := make([]int, 4)

	bFeats := make([]int, 4)

	cFeats := make([]int, 4)

	dFeats := make([]int, 4)

	// This is the ultimate feature vector
	features := make([]int, 196)
	featPtr := 0

	// Slide a 16x16px window by 8px at a time
	for y := 0; y <= 48; y += 8 {
		for x := 0; x <= 48; x += 8 {

			// x and y now define the coordinates of the upper left corner of our window

			// blank out all feature arrays for this subwindow
			for i := range aFeats {
				aFeats[i] = 0
				bFeats[i] = 0
				cFeats[i] = 0
				dFeats[i] = 0
			}

			// count features

			// vary our feature section based on what subimage we're counting
			var secFeats *[]int
			for yp := 0; yp < 16; yp++ {
				for xp := 0; xp < 16; xp++ {
					// select our range and set our current section pointer
					// TODO:	better way to do section classification?
					// 			more elegant loop?
					//			remove the pointer to section and select another way
					switch {
					case inARange(xp, yp):
						secFeats = &aFeats
					case inBRange(xp, yp):
						secFeats = &bFeats
					case inCRange(xp, yp):
						secFeats = &cFeats
					default:
						secFeats = &dFeats
					}

					feature := dirmat[y+yp][x+xp]

					// complex features reduce to two simple features
					// TODO: come up with a more elegant reduction
					switch feature {
					case 0, 1, 2, 3:
						(*secFeats)[feature]++
					case 4, 5:
						(*secFeats)[0]++
						(*secFeats)[3]++
					case 6, 7:
						(*secFeats)[0]++
						(*secFeats)[2]++
					case 8, 9:
						(*secFeats)[1]++
						(*secFeats)[3]++
					case 10, 11:
						(*secFeats)[1]++
						(*secFeats)[2]++
					}
				}
			}

			// Now that we've populated our section features, weight them and add them to our ultimate vector
			for i := range aFeats {
				features[featPtr] = 4*aFeats[i] + 3*bFeats[i] + 2*cFeats[i] + dFeats[i]
				featPtr++
			}
		}
	}

	return features
}

// euclideanDistance measures the distance in Euclidean space between two vectors of like dimension.
func euclideanDistance(vec1, vec2 []int) float64 {
	dist := 0.0

	for i := range vec1 {
		sub := float64(vec1[i] - vec2[i])
		dist += sub * sub
	}

	return math.Sqrt(dist)
}

// KanjiClass determines which kanji a given 196-dimensional vector represents.
// The function takes a kanji vector and a map of kanji names to slices of kanji vectors.
// Each kanji may have more than one training vector. The distance is averaged over all of them.
// The kanji name with the lowest average Euclidean distance to the input vector is returned.
func KanjiClass(kvec []int, vecdb map[string][][]int) string {
	bestAvg := -1.0
	var bestKanji string
	for kanji, vecs := range vecdb {
		avg := 0.0
		denom := 0.0
		for i := 0; i < len(vecs); i++ {
			avg += euclideanDistance(kvec, vecs[i])
			denom += 1.0
		}
		avg /= denom

		if bestAvg == -1 || avg < bestAvg {
			bestAvg = avg
			bestKanji = kanji
		}
	}

	return bestKanji
}

