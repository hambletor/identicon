package main

import (
	"image/color"
	"log"

	"github.com/hambletor/identicon"
)

func main() {

	// Create a simple icon
	simple, err := identicon.New("simple")
	if err != nil {
		log.Printf("error creating icon %v\n", err)
	}

	// save the icon to a file named simple.png
	err = simple.SavePNG()
	if err != nil {
		log.Printf("error saving simple.png, %v\n", err)
	}

	// Create a custom icon
	custom, err := identicon.New("custom",
		identicon.WithBackgroundColor(color.Black), // change the background color
		identicon.WithForegroundColor(color.White), // change the foreground color
		identicon.WithSize(7),                      // change the blocks per row & column block x block
		identicon.WithPixels(350))                  // change the size of icon's side pixels x pixels square

	// save the custom icon to a file named custom.jpeg
	err = custom.SaveJPG()
	if err != nil {
		log.Printf("erro saving custom.jpeg, %v\n", err)
	}
}
