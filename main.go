package main

import (
	"flag"
	"golang.org/x/image/font/sfnt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ppem    = 8
	width   = 8
	height  = 8
	originX = 1
	originY = 5
)

func main() {
	fontFile := flag.String("font", "", "True/Open Type font file.")
	trim := flag.Bool("trim", true, "trim bitmap data of redundant empty values")
	hash := flag.Bool("hash", false, "Replace non-displaying characters with hash blocks")
	credit := flag.String("credit", "", "Font credit to add in the comments section")
	pkgName := flag.String("pkg", "main", "Package name to use in font file")
	fontName := flag.String("name", "", "Font name")
	output := flag.String("o", "", "full filepath (include .go extension) for the output file")
	//mono := flag.Bool("mono", false, "create mono-space fonts 6 bits wide")
	flag.Parse()

	if *fontFile == "" || *output == "" {
		log.Fatalf("Please specify both font file (-font) and output paths (-o)")
	}
	if *fontName == "" {
		fname := filepath.Base(*fontFile)
		fname = strings.Replace(fname, filepath.Ext(fname), "", 1)
		fontName = &fname
	}

	fontData, err := os.ReadFile(*fontFile)
	if err != nil {
		log.Fatalf("Read font file: %v", fontData)
	}

	f, err := sfnt.Parse(fontData)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}

	fontPkg, err := rasterize(f, RasterizeOptions{
		Trim:      *trim,
		HashBlock: *hash,
	})
	if err != nil {
		log.Fatalf("rasterize: %v", err)
	}

	/*
	  top := len(fontPkg.Glyphs)
	  for i:=0; i< top; i++ {
	    g := fontPkg.Glyphs[i]
	    bits:=make([]string, len(g.Bitmaps))
	    for j:=0; j<len(g.Bitmaps); j++ {
	      bits[j] = fmt.Sprintf("0x%02x",g.Bitmaps[j])
	    }
	    fmt.Printf("%d: %s - {%s}\n", g.Rune, string(g.Rune), strings.Join(bits,","))
	  }
	*/
	if err := CreateFontFile(FontFileOptions{
		filepath:    *output,
		fontname:    *fontName,
		packageName: *pkgName,
		credit:      *credit,
	}, fontPkg); err != nil {
		log.Fatal(err)
	}

}
