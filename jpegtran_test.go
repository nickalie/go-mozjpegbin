package mozjpegbin

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestJpegTranReader(t *testing.T) {
	c := NewJpegTran()
	f, err := os.Open("source.jpg")
	assert.Nil(t, err)
	c.Input(f)
	c.OutputFile("target.jpg")
	err = c.Run()
	assert.Nil(t, err)
	validateJpg(t)
}

func TestJpegTranFile(t *testing.T) {
	c := NewJpegTran()
	c.InputFile("source.jpg")
	c.OutputFile("target.jpg")
	err := c.Run()
	assert.Nil(t, err)
	validateJpg(t)
}

func TestJpegTranCrop(t *testing.T) {
	c := NewJpegTran()
	c.Crop(100, 100, 100, 100)
	c.InputFile("source.jpg")
	c.OutputFile("target.jpg")
	err := c.Run()
	assert.Nil(t, err)
}

func TestJpegTranWriter(t *testing.T) {
	f, err := os.Create("target.jpg")
	assert.Nil(t, err)
	defer f.Close() //nolint:errcheck // backup close; validateJpg fails if the file is incomplete

	c := NewJpegTran()
	c.InputFile("source.jpg")
	c.Output(f)
	err = c.Run()
	assert.Nil(t, err)
	f.Close() //nolint:errcheck // validateJpg re-reads the file and would fail on a bad write
	validateJpg(t)
}

func TestJpegTranVersion(t *testing.T) {
	v, err := NewJpegTran().Version()
	assert.Nil(t, err)
	assert.NotZero(t, v)
}
