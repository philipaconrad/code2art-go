// code2art program - Turns text files into black-and-white pixel art!
// Copyright (c) 2022, Philip Conrad. All rights reserved.
package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

func main() {
	// Basic usage printout to ensure we get an argument.
	if argLength := len(os.Args[1:]); argLength < 1 {
		log.Println("Need at least one filename argument.")
		os.Exit(1)
	}

	// Naively open a file using the first CLI arg as the filename.
	fileReader, error := os.Open(os.Args[1])
	if error != nil {
		log.Fatal(error)
	}
	defer fileReader.Close()

	// Read in all of the text into an array, and get length of longest line.
	var lines []string
	longest_line := 0
	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		//fmt.Println(len(scanner.Text()))
		line := scanner.Text()
		lines = append(lines, line)
		if len(line) > longest_line {
			longest_line = len(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Create a new image, and fill in black pixels for text, and
	// white pixels for whitespace and filler.
	width, height := longest_line, len(lines)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	white := color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	black := color.NRGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}
	for y := 0; y < height; y++ {
		line := lines[y]
		for x := 0; x < width; x++ {
			if x < len(line) {
				if unicode.IsSpace(rune(line[x])) {
					img.Set(x, y, white)
				} else {
					img.Set(x, y, black)
				}
			} else {
				img.Set(x, y, white)
			}
		}
	}

	// Open PNG file in current directory, and ensure it gets closed later.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	outFilename := filepath.Base(os.Args[1])
	f, err := os.Create(dir + "/" + outFilename + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Dump pixels to file.
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}

	// Print debug helper message.
	fmt.Printf("Wrote (%d, %d) PNG image to %s\n", width, height, dir+"/"+outFilename+".png")
}
