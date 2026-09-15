package colorspace

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// StreamConvert reads one color per line from r (hex, rgb(), cmyk(), or a
// named color, see Parse) and writes each color, converted to the target
// space, as one line to w.
//
// It holds at most one line in memory at a time (via bufio.Scanner) and
// flushes the output writer once at the end, so a palette file with
// millions of lines costs the same handful of kilobytes as one with ten.
// Callers that need input of unbounded line length should wrap r
// themselves; StreamConvert caps individual lines at 64KiB, which is far
// beyond any real color token, as a guard against unbounded buffering on
// malformed input rather than a limit meant to be tuned.
func StreamConvert(r io.Reader, w io.Writer, target string) error {
	convert, err := converterFor(target)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 4096), 64*1024)

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		rgb, err := Parse(text)
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNo, err)
		}

		if _, err := fmt.Fprintln(bw, convert(rgb)); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("colorspace: reading input: %w", err)
	}

	return bw.Flush()
}

func converterFor(target string) (func(RGB) string, error) {
	switch target {
	case "lab":
		return func(c RGB) string {
			l := c.ToLab()
			return fmt.Sprintf("L=%.2f a=%.2f b=%.2f", l.L, l.A, l.B)
		}, nil
	case "hsl":
		return func(c RGB) string {
			h := c.ToHSL()
			return fmt.Sprintf("H=%.1f S=%.3f L=%.3f", h.H, h.S, h.L)
		}, nil
	case "xyz":
		return func(c RGB) string {
			x := c.ToXYZ()
			return fmt.Sprintf("X=%.4f Y=%.4f Z=%.4f", x.X, x.Y, x.Z)
		}, nil
	case "cmyk":
		return func(c RGB) string {
			k := c.ToCMYK()
			return fmt.Sprintf("C=%.3f M=%.3f Y=%.3f K=%.3f", k.C, k.M, k.Y, k.K)
		}, nil
	default:
		return nil, fmt.Errorf("colorspace: unknown target space %q (want lab, hsl, xyz, or cmyk)", target)
	}
}
