package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"
	"tinygo.org/x/tinyfont"
)

type FontFileOptions struct {
	filepath, fontname, packageName, credit string
}

//go:embed template.txt
var tmpl string

var _ embed.FS

func CreateFontFile(opts FontFileOptions, font *tinyfont.Font) error {
	if opts.packageName == "" {
		opts.packageName = "main"
	}

	funcMap := template.FuncMap{
		"title": strings.Title,
		"str":   runeToString,
		"bstr":  bitmapsToHexArrayString,
	}

	t, err := template.New("fontFile").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	fle, err := os.OpenFile(opts.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fle.Close()
	return t.Execute(fle, map[string]interface{}{
		"PackageName": opts.packageName,
		"FontName":    opts.fontname,
		"Glyphs":      font.Glyphs,
		"YAdvance":    font.YAdvance,
		"Credit":      opts.credit,
	})

}

func runeToString(r rune) string { return string(r) }

func bitmapsToHexArrayString(bitmaps []byte) string {
	bits := make([]string, len(bitmaps))
	for j := 0; j < len(bitmaps); j++ {
		bits[j] = fmt.Sprintf("0x%02x", bitmaps[j])
	}
	return strings.Join(bits, ",")
}
