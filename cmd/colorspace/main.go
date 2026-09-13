// Command colorspace converts sRGB hex colors, one per line, from stdin
// to another color space on stdout. It exists to keep the CLI a thin
// wrapper: all the conversion logic lives in the colorspace package.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spencer-i221/colorspace-cli/colorspace"
)

func main() {
	target := flag.String("to", "lab", "target color space: lab, hsl, or xyz")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "colorspace converts sRGB hex colors read one per line from stdin.\n\n")
		fmt.Fprintf(os.Stderr, "usage:\n  colorspace -to lab < colors.txt\n  cat colors.txt | colorspace -to hsl\n\n")
		fmt.Fprintf(os.Stderr, "each input line is a hex color like \"#ff8800\"; blank lines are skipped.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := colorspace.StreamConvert(os.Stdin, os.Stdout, *target); err != nil {
		fmt.Fprintf(os.Stderr, "colorspace: %v\n", err)
		os.Exit(1)
	}
}
