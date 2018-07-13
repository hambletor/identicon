package identicon

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

const (
	defaultPixels = 250
	//MinPixels defines the smallest icon size in pixels
	MinPixels = 100
	//MaxPixels defines the largest icon size in pixels
	MaxPixels = 500

	defaultSize = 5
	//MinSize is the smallest n can be in an n x n pattern
	MinSize = 5
	//MaxSize is the largest n can be in an n x n pattern
	MaxSize = 20
)

//Icon holds all of the info to build a mirrored block icon
type Icon struct {
	size       int //n x n blocks
	pixels     int // m x m pixels where a block is m/n pixels high and wide
	file       string
	checksum   []byte
	grid       []bool
	img        image.Image
	foreground color.Color
	background color.Color
}

//Option is a functional option used to apply new options (Thanks Dave & Fransesc)
type Option func(*Icon) error

//New creates a new Icon image given file name and Options
func New(name string, options ...Option) (*Icon, error) {
	//name needs to be at least one rune long
	if len(name) == 0 {
		return nil, fmt.Errorf("invalid icon name entered: %s", name)
	}
	i := &Icon{
		file:       name,
		size:       defaultSize,
		pixels:     defaultPixels,
		background: color.White,
	}

	// create a checksum from the name of the file as the seed for the color and grid
	i.checksum = createChecksum(i.file)

	// set the background color using the first three bytes of the checksum
	i.foreground = color.RGBA{
		R: i.checksum[0],
		G: i.checksum[1],
		B: i.checksum[2],
		A: 255,
	}

	errors := make([]error, 0)
	for _, option := range options {
		err := option(i)
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		msg := "\nOption errors:"
		for _, e := range errors {
			msg = fmt.Sprintf("%s\n%s", msg, e)
		}
		return nil, fmt.Errorf("Config error: %s", msg)
	}

	// create the n x n grid as single slice of booleans (if the grid element is true, print a block foreground color)
	i.grid = make([]bool, i.size*i.size)
	createPatternGrid(i)
	i.draw()
	return i, nil
}

//WithPixels set the number of pixels per side of the icon
func WithPixels(p int) Option {
	return func(i *Icon) error {
		if p < MinPixels {
			return fmt.Errorf("pixel length needs to be greater than %d", MinPixels)
		}
		if p > MaxPixels {
			return fmt.Errorf("pixel length needs to be less than %d", MaxPixels)
		}
		i.pixels = int(p)
		return nil
	}
}

//WithSize sets the number of blocks per column and row
func WithSize(s int) Option {
	return func(i *Icon) error {
		if s < MinSize {
			return fmt.Errorf("grid size can not be less than %d", MinSize)
		}
		if s > MaxSize {
			return fmt.Errorf("grid size can not exceed max of %d", MaxSize)
		}
		i.size = int(s)
		return nil
	}
}

//WithBackgroundColor sets the background color
func WithBackgroundColor(c color.Color) Option {
	return func(i *Icon) error {
		if c == nil {
			return fmt.Errorf("can not set background to nil color, please provide a valid color")
		}
		i.background = c
		return nil
	}
}

//WithForegroundColor sets the foreground color
func WithForegroundColor(c color.Color) Option {
	return func(i *Icon) error {
		if c == nil {
			return fmt.Errorf("can not set foreground color to nil, please provide a valid color")
		}
		i.foreground = c
		return nil
	}
}

// allows for changing of the checksum to include more elements
func createChecksum(input ...string) []byte {
	ck := md5.New()
	for _, s := range input {
		io.WriteString(ck, s)
	}
	sum := ck.Sum(nil)
	//need to ensure that the length of the checksum is greater than three
	return sum
}

func createPatternGrid(i *Icon) {
	// for n x n given an checksum array of length l build a mirror matrix
	// if l / (n/2) > n we have enough data to just use mirrored slices of odd/even to create grid
	chunk := (i.size / 2)
	if i.size%2 == 1 {
		chunk++
	}
	// build the data (minus mirrored data) to establish a pattern
	var data []byte
	done := false
	j := 0
	for !done {
		data = append(data, i.checksum[j])
		if j == len(i.checksum)-1 {
			j = 0
		}
		if len(data) == chunk*i.size {
			done = true
		}
		j++
	}

	// now create each row of the icon with []data created from checksum
	for row := 0; row < i.size; row++ {
		//grab a slice of the grid that represents a row
		r := i.grid[row*i.size : (row*i.size)+i.size]
		//get a slice that represents the pattern
		c := data[row*chunk : (row*chunk)+chunk]

		//create the mirrored pattern from left to center then right to center
		for j, m := 0, 0; j < len(r); j++ {
			//from left to ceter
			if j < len(c) {
				r[j] = int(c[j])%2 == 0 //only even bytes from checksum data are marked true and printed
				continue
			}
			// mirror from the right to center
			r[len(r)-1-m] = int(c[m])%2 == 0 //only even bytes from checksum data are marked true and printed
			m++
		}
	}
}

func (i *Icon) draw() error {
	//go through grid to create the image
	base := image.NewRGBA(image.Rect(0, 0, i.pixels, i.pixels)) // build a base image of pixels squared

	// draw base icon square with background fill of Icon
	draw.Draw(base, base.Bounds(), &image.Uniform{i.background}, image.ZP, draw.Src)

	// start drawing pattern using the boolean pattern grid (if true draw a block) using Icon foreground
	length := i.pixels / i.size
	for j := 0; j < len(i.grid); j++ {
		if i.grid[j] {
			x1 := (j % i.size) * length
			y1 := (j / i.size) * length
			x2 := x1 + length
			y2 := y1 + length
			block := image.Rect(x1, y1, x2, y2)
			draw.Draw(base, block, &image.Uniform{i.foreground}, image.ZP, draw.Src)
		}
	}

	i.img = base
	return nil
}

//Pattern display's the pattern of the image
func (i *Icon) Pattern() string {
	return i.pattern()
}

//SavePNG saves the icon to png format
func (i *Icon) SavePNG() error {
	return saveFile(&i.img, i.file+".png", false)
}

//SaveJPG saves the icon to jpg format
func (i *Icon) SaveJPG() error {
	return saveFile(&i.img, i.file+".jpeg", true)
}

func saveFile(i *image.Image, name string, jpg bool) error {
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("unable to create file %s", name)
	}
	if jpg {
		err = jpeg.Encode(file, *i, nil)
		if err != nil {
			return fmt.Errorf("issue saving jpeg %v", err)
		}
	}
	err = png.Encode(file, *i)
	if err != nil {
		return fmt.Errorf("issue saving png %v", err)
	}
	return nil
}

//String satisfies the Stringer interface
func (i Icon) String() string {
	return fmt.Sprintf("Icon:\nfile name: %s\nsize in pixels %d x %d\nsize in blocks %d x %d\nforeground %v\nbackground %v\n%s",
		i.file, i.pixels, i.pixels, i.size, i.size, i.foreground, i.background, i.pattern())
}

func (i *Icon) pattern() string {
	var s string
	s = "Pattern:"
	for g := 0; g < len(i.grid); g++ {
		if g%i.size == 0 {
			s = s + "\n"
		}
		if i.grid[g] {
			s = s + "*"
		} else {
			s = s + "."
		}
	}
	return s
}
