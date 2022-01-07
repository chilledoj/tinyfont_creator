# ARRGGGGGHHHH!

Whoops. It doesn't work! I misunderstood how the data in the bitmaps is structured!! **Do not use**.

**TODO**: FIX - Change mapping to horizontal instead of vertical rendering.

# TinyFont Creator

Creates a `.go` file with bitmap font images to be used by [TinyFont](https://github.com/tinygo-org/tinyfont) from a truetype font file.

## Usage
`tfCreator -font ./dot_matrix/DOTMATRI.TTF -hash=true -credit "&copy; Svein KÂre Gunnarson" -o ./dotMatrix.go`

### Options:

| Flag | Type | Default | Description |
| ---- | :--: | :-----: | ----------- |
| `-font` | string | "" | Full filepath to font file (required) |
| `-o` | string | "" | full filepath to output font file - ensure you include the .go extension |
| `-trim` | bool | true | if true, will trim excess empty space from the vectors (zeroes) to ensure a variable width font. |
| `-hash` | bool | false | if true, will render a 4 bit width hash block instead |
| `-pkg` | string | main | Package name to use in the generated go file. |
| `-name` | string | _filename_ | Name of font used in variable declaration. By default it will use the filename of the provided TTF. |
| `-credit` | string | "" | Text to add to the top of the go file as a comment to provide a credit to the font creator. |

