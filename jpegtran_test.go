package mozjpegbin

import (
	"os"
	"github.com/stretchr/testify/assert"
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
	c.Crop(500, 500, 100, 100)
	c.InputFile("source.jpg")
	c.OutputFile("target.jpg")
	err := c.Run()
	assert.Nil(t, err)
}

func TestJpegTranWriter(t *testing.T) {
	f, err := os.Create("target.jpg")
	assert.Nil(t, err)
	defer f.Close()

	c := NewJpegTran()
	c.InputFile("source.jpg")
	c.Output(f)
	err = c.Run()
	assert.Nil(t, err)
	f.Close()
	validateJpg(t)
}

func TestJpegTranVersion(t *testing.T) {
	v, err := NewJpegTran().Version()
	assert.Nil(t, err)
	assert.NotZero(t, v)
}
