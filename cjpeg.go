package mozjpegbin

import (
	"github.com/nickalie/go-binwrapper"
	"image"
	"io"
	"fmt"
	"errors"
)

type CJpeg struct {
	*binwrapper.BinWrapper
	inputFile  string
	inputImage image.Image
	input      io.Reader
	outputFile string
	output     io.Writer
	quality    int
	optimize   bool
}

func NewCJpeg() *CJpeg {
	bin := &CJpeg{
		BinWrapper: createBinWrapper(),
		quality:    -1,
	}
	bin.ExecPath("cjpeg")

	return bin
}

//Sets image file to convert
//Input or InputImage called before will be ignored
func (c *CJpeg) InputFile(file string) *CJpeg {
	c.input = nil
	c.inputImage = nil
	c.inputFile = file
	return c
}

//Sets reader to convert
//InputFile or InputImage called before will be ignored
func (c *CJpeg) Input(reader io.Reader) *CJpeg {
	c.inputFile = ""
	c.inputImage = nil
	c.input = reader
	return c
}

//Sets image to convert
//InputFile or Input called before will be ignored
func (c *CJpeg) InputImage(img image.Image) *CJpeg {
	c.inputFile = ""
	c.input = nil
	c.inputImage = img
	return c
}

//Specify the name of the output JPEG file
//Output called before will be ignored
func (c *CJpeg) OutputFile(file string) *CJpeg {
	c.output = nil
	c.outputFile = file
	return c
}

//Specify writer to write JPEG file content
//OutputFile called before will be ignored
func (c *CJpeg) Output(writer io.Writer) *CJpeg {
	c.outputFile = ""
	c.output = writer
	return c
}

//Specify the compression factor for RGB channels between 0 and 100. The default is 75.
//
//A small factor produces a smaller file with lower quality. Best quality is achieved by using a value of 100.
func (c *CJpeg) Quality(quality uint) *CJpeg {
	if quality > 100 {
		quality = 100
	}

	c.quality = int(quality)
	return c
}

//Perform optimization of entropy encoding parameters.
//Without this, default encoding parameters are used.
//Optimize usually makes the JPEG file a little smaller, but cjpeg runs somewhat slower and needs much more memory.
//Image quality and speed of decompression are unaffected by Optimize.
func (c *CJpeg) Optimize(optimize bool) *CJpeg {
	c.optimize = optimize
	return c
}

//Runs encoder
func (c *CJpeg) Run() error {
	defer c.BinWrapper.Reset()

	if c.quality > -1 {
		c.Arg("-quality", fmt.Sprintf("%d", c.quality))
	}

	if c.optimize {
		c.Arg("-optimize")
	}

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

func (c *CJpeg) Version() (string, error)  {
	return version(c.BinWrapper)
}

//Resets all parameters to default values
func (c *CJpeg) Reset() *CJpeg {
	c.quality = -1
	c.optimize = false
	return c
}

func (c *CJpeg) setInput() error {
	if c.input != nil {
		c.StdIn(c.input)
	} else if c.inputImage != nil {
		r, err := createReaderFromImage(c.inputImage)

		if err != nil {
			return err
		}

		c.StdIn(r)
	} else if c.inputFile != "" {
		c.Arg(c.inputFile)
	} else {
		return errors.New("Undefined input")
	}

	return nil
}

func (c *CJpeg) getOutput() (string, error) {
	if c.output != nil {
		return "", nil
	} else if c.outputFile != "" {
		return c.outputFile, nil
	} else {
		return "", errors.New("Undefined output")
	}
}
