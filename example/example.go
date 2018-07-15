package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hambletor/identicon"
)

func main() {

	// Create a simple icon
	simple, err := identicon.New("simple")
	if err != nil {
		log.Printf("error creating icon %v\n", err)
		return
	}

	// save the icon to a file named simple.png
	err = simple.SavePNG()
	if err != nil {
		log.Printf("error saving simple.png, %v\n", err)
	}
	fmt.Printf("%v\n", simple.Pattern())

	fmt.Print("\n--------------\n\n")
	fg := color.RGBA{R: 122, G: 16, B: 21, A: 255}
	// Create a custom icon
	custom, err := identicon.New("custom",
		// identicon.WithBackgroundColor(bg), // change the background color
		identicon.WithForegroundColor(fg),       // change the foreground color
		identicon.WithComplementaryBackground(), // makes the background a complimentary color
		identicon.WithSize(13),                  // change the blocks per row & column block x block
		identicon.WithPixels(500))               // change the size of icon's side pixels x pixels square

	if err != nil {
		log.Printf("error creating custom identicon %v", err)
		return
	}
	// save the custom icon to a file named custom.jpeg
	err = custom.SaveJPG()
	if err != nil {
		log.Printf("erro saving custom.jpeg, %v\n", err)
	}

	fmt.Printf("%v", custom)

}
