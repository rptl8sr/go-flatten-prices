package processors

import (
	"bufio"
	"encoding/csv"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"
)

type Tag struct {
	Code  int
	Font  string
	Size  int
	Color color.RGBA
	X     int
	Y     int
	Align string
	Price *int
}

func hexToRGBA(hex string) (c color.RGBA, err error) {
	hex = strings.TrimPrefix(hex, "#")

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return c, err
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return c, err
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return c, err
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff}, nil
}

func ReadTags(filename string) ([]*Tag, error) {
	f, e := os.Open(filename)
	if e != nil {
		return nil, e
	}
	defer func() { _ = f.Close() }()

	reader := csv.NewReader(bufio.NewReader(f))
	reader.Comma = ';'
	reader.LazyQuotes = true

	var tags []*Tag
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		code, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		size, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, err
		}

		fontColor, err := hexToRGBA(record[3])
		if err != nil {
			return nil, err
		}

		x, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, err
		}

		y, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, err
		}

		tags = append(tags, &Tag{
			Code:  code,
			Font:  record[1],
			Size:  size,
			Color: fontColor,
			X:     x,
			Y:     y,
			Align: record[6],
		})
	}

	return tags, nil
}
