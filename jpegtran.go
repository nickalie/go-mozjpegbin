package mozjpegbin

import (
	"errors"
	"fmt"
	"github.com/nickalie/go-binwrapper"
	"io"
)

type cropInfo struct {
	x      int
	y      int
	width  int
	height int
}

// JpegTran wraps jpegtran tool from mozjpeg
type JpegTran struct {
	*binwrapper.BinWrapper
	optimize    bool
	progressive bool
	crop        *cropInfo
	inputFile   string
	input       io.Reader
	outputFile  string
	output      io.Writer
	copy        string
}

// NewJpegTran creates new JpegTran instance
func NewJpegTran() *JpegTran {
	bin := &JpegTran{
		BinWrapper: createBinWrapper(),
		copy:       "none",
		optimize:   true,
	}

	bin.ExecPath("jpegtran")
	return bin
}

// Optimize perform optimization of entropy encoding parameters
func (c *JpegTran) Optimize(optimize bool) *JpegTran {
	c.optimize = optimize
	return c
}

// Progressive create progressive JPEG file
func (c *JpegTran) Progressive(progressive bool) *JpegTran {
	c.progressive = progressive
	return c
}

// Crop to a rectangular region of width and height, starting at point x,y
func (c *JpegTran) Crop(x, y, width, height int) *JpegTran {
	c.crop = &cropInfo{x, y, width, height}
	return c
}

// InputFile sets image file to convert.
// Input or InputImage called before will be ignored.
func (c *JpegTran) InputFile(file string) *JpegTran {
	c.input = nil
	c.inputFile = file
	return c
}

// Input sets reader to convert.
// InputFile or InputImage called before will be ignored.
func (c *JpegTran) Input(reader io.Reader) *JpegTran {
	c.inputFile = ""
	c.input = reader
	return c
}

// OutputFile specify the name of the output jpeg file.
// Output called before will be ignored.
func (c *JpegTran) OutputFile(file string) *JpegTran {
	c.output = nil
	c.outputFile = file
	return c
}

// Output specify writer to write jpeg file content.
// OutputFile called before will be ignored.
func (c *JpegTran) Output(writer io.Writer) *JpegTran {
	c.outputFile = ""
	c.output = writer
	return c
}

// CopyNone copy no extra markers from source file. This setting suppresses all comments and other metadata in the source file
func (c *JpegTran) CopyNone() *JpegTran {
	c.copy = "none"
	return c
}

// CopyComments copy only comment markers.  This setting copies comments from the source file but discards any other metadata.
func (c *JpegTran) CopyComments() *JpegTran {
	c.copy = "comments"
	return c
}

// CopyAll copy all extra markers. This setting preserves miscellaneous markers found in the source file, such as JFIF thumbnails, Exif data, and Photoshop settings. In some files, these extra markers can be sizable. Note that this option will copy thumbnails as-is; they will not be transformed.
func (c *JpegTran) CopyAll() *JpegTran {
	c.copy = "all"
	return c
}

// Run starts jpegtran with specified parameters.
func (c *JpegTran) Run() error {
	defer c.BinWrapper.Reset()

	if c.optimize {
		c.Arg("-optimize")
	}

	if c.progressive {
		c.Arg("-progressive")
	}

	if c.crop != nil {
		c.Arg("-crop", fmt.Sprintf("%dx%d+%d+%d", c.crop.width, c.crop.height, c.crop.x, c.crop.y))
	}

	c.Arg("-copy", c.copy)

	output, err := c.getOutput()

	if err != nil {
		return err
	}

	if output != "" {
		c.Arg("-outfile", output)
	}

	err = c.setInput()

	if err != nil {
		return err
	}

	if c.output != nil {
		c.SetStdOut(c.output)
	}

	err = c.BinWrapper.Run()

	if err != nil {
		return errors.New(err.Error() + ". " + string(c.StdErr()))
	}

	return nil
}

// Version returns jpegtran version.
func (c *JpegTran) Version() (string, error) {
	return version(c.BinWrapper)
}

// Reset resets all parameters to default values
func (c *JpegTran) Reset() *JpegTran {
	c.optimize = true
	c.progressive = false
	c.copy = "none"
	c.crop = nil
	return c
}

func (c *JpegTran) setInput() error {
	if c.input != nil {
		c.StdIn(c.input)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *JpegTran) getOutput() (string, error) {
	if c.output != nil {
		return "", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
