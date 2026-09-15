package colorspace

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse parses a color given in any of the formats this package accepts:
// a hex string ("#ff8800" or "ff8800"), an "rgb(r, g, b)" triple with each
// channel in [0,255], a "cmyk(c, m, y, k)" quadruple with each channel a
// percentage in [0,100], or a CSS/X11 named color such as "cornflowerblue"
// (case-insensitive).
func Parse(s string) (RGB, error) {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)

	switch {
	case strings.HasPrefix(lower, "rgb(") && strings.HasSuffix(lower, ")"):
		return parseRGBFunc(s[len("rgb(") : len(s)-1])
	case strings.HasPrefix(lower, "cmyk(") && strings.HasSuffix(lower, ")"):
		return parseCMYKFunc(s[len("cmyk(") : len(s)-1])
	}

	if hex, ok := namedColors[lower]; ok {
		return ParseHex(hex)
	}

	return ParseHex(s)
}

func parseRGBFunc(body string) (RGB, error) {
	parts := strings.Split(body, ",")
	if len(parts) != 3 {
		return RGB{}, fmt.Errorf("colorspace: rgb() wants 3 comma-separated channels, got %d", len(parts))
	}
	chans := make([]float64, 3)
	for i, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return RGB{}, fmt.Errorf("colorspace: invalid rgb() channel %q: %w", p, err)
		}
		if v < 0 || v > 255 {
			return RGB{}, fmt.Errorf("colorspace: rgb() channel %d out of range [0,255]", v)
		}
		chans[i] = float64(v) / 255
	}
	return RGB{R: chans[0], G: chans[1], B: chans[2]}, nil
}

func parseCMYKFunc(body string) (RGB, error) {
	parts := strings.Split(body, ",")
	if len(parts) != 4 {
		return RGB{}, fmt.Errorf("colorspace: cmyk() wants 4 comma-separated channels, got %d", len(parts))
	}
	chans := make([]float64, 4)
	for i, p := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return RGB{}, fmt.Errorf("colorspace: invalid cmyk() channel %q: %w", p, err)
		}
		if v < 0 || v > 100 {
			return RGB{}, fmt.Errorf("colorspace: cmyk() channel %v out of range [0,100]", v)
		}
		chans[i] = v / 100
	}
	c := CMYK{C: chans[0], M: chans[1], Y: chans[2], K: chans[3]}
	return c.ToRGB(), nil
}

// namedColors maps lowercase CSS/X11 color keywords to their hex value.
var namedColors = map[string]string{
	"aliceblue":            "f0f8ff",
	"antiquewhite":         "faebd7",
	"aqua":                 "00ffff",
	"aquamarine":           "7fffd4",
	"azure":                "f0ffff",
	"beige":                "f5f5dc",
	"bisque":               "ffe4c4",
	"black":                "000000",
	"blanchedalmond":       "ffebcd",
	"blue":                 "0000ff",
	"blueviolet":           "8a2be2",
	"brown":                "a52a2a",
	"burlywood":            "deb887",
	"cadetblue":            "5f9ea0",
	"chartreuse":           "7fff00",
	"chocolate":            "d2691e",
	"coral":                "ff7f50",
	"cornflowerblue":       "6495ed",
	"cornsilk":             "fff8dc",
	"crimson":              "dc143c",
	"cyan":                 "00ffff",
	"darkblue":             "00008b",
	"darkcyan":             "008b8b",
	"darkgoldenrod":        "b8860b",
	"darkgray":             "a9a9a9",
	"darkgreen":            "006400",
	"darkgrey":             "a9a9a9",
	"darkkhaki":            "bdb76b",
	"darkmagenta":          "8b008b",
	"darkolivegreen":       "556b2f",
	"darkorange":           "ff8c00",
	"darkorchid":           "9932cc",
	"darkred":              "8b0000",
	"darksalmon":           "e9967a",
	"darkseagreen":         "8fbc8f",
	"darkslateblue":        "483d8b",
	"darkslategray":        "2f4f4f",
	"darkslategrey":        "2f4f4f",
	"darkturquoise":        "00ced1",
	"darkviolet":           "9400d3",
	"deeppink":             "ff1493",
	"deepskyblue":          "00bfff",
	"dimgray":              "696969",
	"dimgrey":              "696969",
	"dodgerblue":           "1e90ff",
	"firebrick":            "b22222",
	"floralwhite":          "fffaf0",
	"forestgreen":          "228b22",
	"fuchsia":              "ff00ff",
	"gainsboro":            "dcdcdc",
	"ghostwhite":           "f8f8ff",
	"gold":                 "ffd700",
	"goldenrod":            "daa520",
	"gray":                 "808080",
	"grey":                 "808080",
	"green":                "008000",
	"greenyellow":          "adff2f",
	"honeydew":             "f0fff0",
	"hotpink":              "ff69b4",
	"indianred":            "cd5c5c",
	"indigo":               "4b0082",
	"ivory":                "fffff0",
	"khaki":                "f0e68c",
	"lavender":             "e6e6fa",
	"lavenderblush":        "fff0f5",
	"lawngreen":            "7cfc00",
	"lemonchiffon":         "fffacd",
	"lightblue":            "add8e6",
	"lightcoral":           "f08080",
	"lightcyan":            "e0ffff",
	"lightgoldenrodyellow": "fafad2",
	"lightgray":            "d3d3d3",
	"lightgreen":           "90ee90",
	"lightgrey":            "d3d3d3",
	"lightpink":            "ffb6c1",
	"lightsalmon":          "ffa07a",
	"lightseagreen":        "20b2aa",
	"lightskyblue":         "87cefa",
	"lightslategray":       "778899",
	"lightslategrey":       "778899",
	"lightsteelblue":       "b0c4de",
	"lightyellow":          "ffffe0",
	"lime":                 "00ff00",
	"limegreen":            "32cd32",
	"linen":                "faf0e6",
	"magenta":              "ff00ff",
	"maroon":               "800000",
	"mediumaquamarine":     "66cdaa",
	"mediumblue":           "0000cd",
	"mediumorchid":         "ba55d3",
	"mediumpurple":         "9370db",
	"mediumseagreen":       "3cb371",
	"mediumslateblue":      "7b68ee",
	"mediumspringgreen":    "00fa9a",
	"mediumturquoise":      "48d1cc",
	"mediumvioletred":      "c71585",
	"midnightblue":         "191970",
	"mintcream":            "f5fffa",
	"mistyrose":            "ffe4e1",
	"moccasin":             "ffe4b5",
	"navajowhite":          "ffdead",
	"navy":                 "000080",
	"oldlace":              "fdf5e6",
	"olive":                "808000",
	"olivedrab":            "6b8e23",
	"orange":               "ffa500",
	"orangered":            "ff4500",
	"orchid":               "da70d6",
	"palegoldenrod":        "eee8aa",
	"palegreen":            "98fb98",
	"paleturquoise":        "afeeee",
	"palevioletred":        "db7093",
	"papayawhip":           "ffefd5",
	"peachpuff":            "ffdab9",
	"peru":                 "cd853f",
	"pink":                 "ffc0cb",
	"plum":                 "dda0dd",
	"powderblue":           "b0e0e6",
	"purple":               "800080",
	"rebeccapurple":        "663399",
	"red":                  "ff0000",
	"rosybrown":            "bc8f8f",
	"royalblue":            "4169e1",
	"saddlebrown":          "8b4513",
	"salmon":               "fa8072",
	"sandybrown":           "f4a460",
	"seagreen":             "2e8b57",
	"seashell":             "fff5ee",
	"sienna":               "a0522d",
	"silver":               "c0c0c0",
	"skyblue":              "87ceeb",
	"slateblue":            "6a5acd",
	"slategray":            "708090",
	"slategrey":            "708090",
	"snow":                 "fffafa",
	"springgreen":          "00ff7f",
	"steelblue":            "4682b4",
	"tan":                  "d2b48c",
	"teal":                 "008080",
	"thistle":              "d8bfd8",
	"tomato":               "ff6347",
	"turquoise":            "40e0d0",
	"violet":               "ee82ee",
	"wheat":                "f5deb3",
	"white":                "ffffff",
	"whitesmoke":           "f5f5f5",
	"yellow":               "ffff00",
	"yellowgreen":          "9acd32",
}
