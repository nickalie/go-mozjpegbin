package mozjpegbin

import (
	"github.com/nickalie/go-binwrapper"
	"io"
	"fmt"
	"errors"
)

type cropInfo struct {
	x      int
	y      int
	width  int
	height int
}

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

func NewJpegTran() *JpegTran {
	bin := &JpegTran{
		BinWrapper: createBinWrapper(),
		copy:       "none",
		optimize:   true,
	}

	bin.ExecPath("jpegtran")
	return bin
}

//Perform optimization of entropy encoding parameters
func (c *JpegTran) Optimize(optimize bool) *JpegTran {
	c.optimize = optimize
	return c
}

//Create progressive JPEG file
func (c *JpegTran) Progressive(progressive bool) *JpegTran {
	c.progressive = progressive
	return c
}

//Crop to a rectangular region of width and height, starting at point x,y
func (c *JpegTran) Crop(x, y, width, height int) *JpegTran {
	c.crop = &cropInfo{x, y, width, height}
	return c
}

//Sets image file to convert
//Input or InputImage called before will be ignored
func (c *JpegTran) InputFile(file string) *JpegTran {
	c.input = nil
	c.inputFile = file
	return c
}

//Sets reader to convert
//InputFile or InputImage called before will be ignored
func (c *JpegTran) Input(reader io.Reader) *JpegTran {
	c.inputFile = ""
	c.input = reader
	return c
}

//Specify the name of the output JPEG file
//Output called before will be ignored
func (c *JpegTran) OutputFile(file string) *JpegTran {
	c.output = nil
	c.outputFile = file
	return c
}

//Specify writer to write JPEG file content
//OutputFile called before will be ignored
func (c *JpegTran) Output(writer io.Writer) *JpegTran {
	c.outputFile = ""
	c.output = writer
	return c
}

//Copy no extra markers from source file. This setting suppresses all comments and other metadata in the source file
func (c *JpegTran) CopyNone() *JpegTran {
	c.copy = "none"
	return c
}

//Copy only comment markers.  This setting copies comments from the source file but discards any other metadata.
func (c *JpegTran) CopyComments() *JpegTran {
	c.copy = "comments"
	return c
}

//Copy all extra markers. This setting preserves miscellaneous markers found in the source file, such as JFIF thumbnails, Exif data, and Photoshop settings. In some files, these extra markers can be sizable. Note that this option will copy thumbnails as-is; they will not be transformed.
func (c *JpegTran) CopyAll() *JpegTran {
	c.copy = "all"
	return c
}

//Run jpegtran
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

func (c *JpegTran) Version() (string, error)  {
	return version(c.BinWrapper)
}

//Resets all parameters to default values
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
