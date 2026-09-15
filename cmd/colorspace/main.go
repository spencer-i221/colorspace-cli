// Command colorspace converts colors, one per line, from stdin to another
// color space on stdout. It exists to keep the CLI a thin wrapper: all the
// conversion logic lives in the colorspace package.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spencer-i221/colorspace-cli/colorspace"
)

func main() {
	target := flag.String("to", "lab", "target color space: lab, hsl, xyz, or cmyk")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "colorspace converts colors read one per line from stdin.\n\n")
		fmt.Fprintf(os.Stderr, "usage:\n  colorspace -to lab < colors.txt\n  cat colors.txt | colorspace -to hsl\n\n")
		fmt.Fprintf(os.Stderr, "each input line is a hex color (\"#ff8800\"), an rgb() or cmyk()\n")
		fmt.Fprintf(os.Stderr, "function (\"rgb(255,136,0)\", \"cmyk(0,47,100,0)\"), or a named color\n")
		fmt.Fprintf(os.Stderr, "(\"cornflowerblue\"); blank lines are skipped.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := colorspace.StreamConvert(os.Stdin, os.Stdout, *target); err != nil {
		fmt.Fprintf(os.Stderr, "colorspace: %v\n", err)
		os.Exit(1)
	}
}
