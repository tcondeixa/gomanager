package color

import (
	"fmt"
	"strconv"
)

const (
	escChar     = "\x1b"
	csi         = escChar + "["
	FG          = 38
	BG          = 48
	resetColors = csi + "m"
)

func HexToRGB(hex string) ([3]uint8, error) {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	// Parse the hex string
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return [3]uint8{}, err
	}

	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return [3]uint8{}, err
	}

	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return [3]uint8{}, err
	}

	return [3]uint8{uint8(r), uint8(g), uint8(b)}, nil
}

type color struct {
	Hex string
	RGB [3]uint8
}

func newColor(hex string) color {
	rgb, err := HexToRGB(hex)
	if err != nil {
		rgb = [3]uint8{255, 255, 255} // default to white on error
	}

	return color{
		Hex: hex,
		RGB: rgb,
	}
}

type Scheme struct {
	noColor bool
	text    color
	header  color
	err     color
}

func NewScheme(noColor bool, text, header, err string) *Scheme {
	return &Scheme{
		noColor: noColor,
		header:  newColor(header),
		text:    newColor(text),
		err:     newColor(err),
	}
}

func (s *Scheme) Text(text string) string {
	if s.noColor {
		return text
	}

	return csi + fmt.Sprintf(
		"%d;2;%d;%d;%dm%s%s",
		FG,
		s.text.RGB[0],
		s.text.RGB[1],
		s.text.RGB[2],
		text,
		resetColors,
	)
}

func (s *Scheme) Header(text string) string {
	if s.noColor {
		return text
	}

	return csi + fmt.Sprintf(
		"%d;2;%d;%d;%dm%s%s",
		FG,
		s.header.RGB[0],
		s.header.RGB[1],
		s.header.RGB[2],
		text,
		resetColors,
	)
}

func (s *Scheme) Err(text string) string {
	if s.noColor {
		return text
	}

	return csi + fmt.Sprintf(
		"%d;2;%d;%d;%dm%s%s",
		FG,
		s.err.RGB[0],
		s.err.RGB[1],
		s.err.RGB[2],
		text,
		resetColors,
	)
}
