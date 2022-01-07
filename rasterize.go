package main

import (
	"errors"
	"fmt"
	"golang.org/x/image/draw"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
	"image"
	"tinygo.org/x/tinyfont"
)

type RasterizeOptions struct {
	Trim      bool
	HashBlock bool
}

var hashBlock = []byte{0b10101010, 0b01010101, 0b10101010, 0b01010101}

func rasterize(f *sfnt.Font, opts RasterizeOptions) (*tinyfont.Font, error) {
	fnt := tinyfont.Font{
		Glyphs:   make([]tinyfont.Glyph, 0),
		YAdvance: 0,
	}
	for i := 0; i < 0x10FFFF; i++ {
		g, err := rasterizeGlyph(f, rune(i))
		if err != nil {
			continue
		}
		if opts.Trim {
			trimBitmaps(g)
		}
		if opts.HashBlock && len(g.Bitmaps) == 0 {
			g.Bitmaps = hashBlock
			g.Width = 4
			g.XAdvance = 5
		}
		fnt.Glyphs = append(fnt.Glyphs, *g)
	}

	return &fnt, nil
}

func rasterizeGlyph(f *sfnt.Font, chr rune) (*tinyfont.Glyph, error) {
	var b sfnt.Buffer
	x, err := f.GlyphIndex(&b, chr)
	if err != nil {
		return nil, err
	}
	if x == 0 {
		return nil, errors.New("GlyphIndex: no glyph index found for the rune 'Ġ'")
	}
	segments, err := f.LoadGlyph(&b, x, fixed.I(ppem), nil)
	if err != nil {
		return nil, fmt.Errorf("LoadGlyph: %w", err)
	}

	// Translate and scale that glyph as we pass it to a vector.Rasterizer.
	r := vector.NewRasterizer(width, height)
	r.DrawOp = draw.Src
	for _, seg := range segments {
		// The divisions by 64 below is because the seg.Args values have type
		// fixed.Int26_6, a 26.6 fixed point number, and 1<<6 == 64.
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			r.MoveTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
			)
		case sfnt.SegmentOpLineTo:
			r.LineTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
			)
		case sfnt.SegmentOpQuadTo:
			r.QuadTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
				originX+float32(seg.Args[1].X)/64,
				originY+float32(seg.Args[1].Y)/64,
			)
		case sfnt.SegmentOpCubeTo:
			r.CubeTo(
				originX+float32(seg.Args[0].X)/64,
				originY+float32(seg.Args[0].Y)/64,
				originX+float32(seg.Args[1].X)/64,
				originY+float32(seg.Args[1].Y)/64,
				originX+float32(seg.Args[2].X)/64,
				originY+float32(seg.Args[2].Y)/64,
			)
		}
	}

	// Finish the rasterization: the conversion from vector graphics (shapes)
	// to raster graphics (pixels).
	dst := image.NewAlpha(image.Rect(0, 0, width, height))
	r.Draw(dst, dst.Bounds(), image.Opaque, image.Point{})

	// Visualize the pixels.

	tg := tinyfont.Glyph{
		Rune:     chr,
		Width:    width,
		Height:   height,
		XAdvance: width + 1,
		XOffset:  0,
		YOffset:  0,
		Bitmaps:  make([]byte, 0),
	}

	// Convert each vertical line into a uint8
	var threshold uint32 = 0x2fff

	for x := 0; x < width; x++ {
		var byt uint8
		for y := 0; y < height; y++ {
			r, g, b, a := dst.At(x, y).RGBA()
			if a > threshold && r > threshold && g > threshold && b > threshold {
				byt |= (1 << y)
			}
		}
		tg.Bitmaps = append(tg.Bitmaps, byt)
	}

	return &tg, nil
}

func trimBitmaps(g *tinyfont.Glyph) {
	// Trim ends
	end := len(g.Bitmaps) - 1
	for i := len(g.Bitmaps) - 1; i >= 0; i-- {
		if g.Bitmaps[i] == 0 {
			end = int(i)
		} else {
			break
		}
	}
	g.Bitmaps = g.Bitmaps[:end]
	g.Width = uint8(end)
	g.XAdvance = g.Width + 1

	// Trim beginning but only if anything left
	if len(g.Bitmaps) == 0 {
		return
	}
	start := 0
	for i := 0; i < len(g.Bitmaps)-1; i++ {
		if g.Bitmaps[i] != 0 {
			start = i
			break
		}
	}
	if start == 0 {
		return
	}
	g.Bitmaps = g.Bitmaps[start:]
	g.Width = uint8(len(g.Bitmaps))
	g.XAdvance = g.Width + 1
	g.XOffset = 1
}
