package identicon

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/fogleman/gg"
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
type Option func(*Icon)

//New creates a new Icon image given file name and Options
func New(name string, options ...Option) (*Icon, error) {
	//name needs to be at least one rune long
	if len(name) == 0 {
		return nil, fmt.Errorf("invalid icon name entered: %s", name)
	}
	i := &Icon{
		file:       name,
		size:       5,
		pixels:     250,
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

	for _, option := range options {
		option(i)
	}
	if i.pixels <= i.size {
		return nil, fmt.Errorf("size:%d needs to be smaller than pixels:%d ", i.size, i.pixels)
	}
	// create the n x n grid as single slice of booleans (if the grid element is true, print a block foreground color)
	i.grid = make([]bool, i.size*i.size)
	createPatternGrid(i)
	i.draw()
	return i, nil
}

//WithPixels set the number of pixels per side of the icon
func WithPixels(p int) Option {
	return func(i *Icon) { i.pixels = p }
}

//WithSize sets the number of blocks per column and row
func WithSize(s int) Option {
	return func(i *Icon) { i.size = s }
}

//WithBackgroundColor sets the background color
func WithBackgroundColor(c color.Color) Option {
	return func(i *Icon) { i.background = c }
}

//WithForegroundColor sets the foreground color
func WithForegroundColor(c color.Color) Option {
	return func(i *Icon) { i.foreground = c }
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
	dc := gg.NewContext(i.pixels, i.pixels) // build a base image of pixels squared

	// draw base icon square with background fill of Icon
	dc.DrawRectangle(0, 0, float64(i.pixels), float64(i.pixels))
	br, bg, bb, _ := i.background.RGBA()
	dc.SetRGB255(int(br), int(bg), int(bb))
	dc.Fill()

	// start drawing pattern using the boolean pattern grid (if true draw a block) using Icon foreground
	fr, fg, fb, _ := i.foreground.RGBA()
	length := i.pixels / i.size
	for j := 0; j < len(i.grid); j++ {
		if i.grid[j] {
			dc.DrawRectangle(float64((j%i.size)*length), float64((j/i.size)*length), float64(length), float64(length))
		}
	}
	dc.SetRGB255(int(fr), int(fg), int(fb))
	dc.Fill()

	i.img = dc.Image()
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
