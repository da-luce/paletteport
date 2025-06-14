package color

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Color represents a color with RGBA components.
type Color struct {
	Alpha float64
	Red   float64
	Green float64
	Blue  float64
}

// NewColor creates a new Color instance.
func NewColor(red, green, blue, alpha float64) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
		Alpha: alpha,
	}
}

// AsString returns the color as a formatted RGBA string.
func (c Color) AsString() string {
	return fmt.Sprintf("RGBA(%.3f, %.3f, %.3f, %.3f)", c.Red, c.Green, c.Blue, c.Alpha)
}

// ToRGB converts the color to an RGB tuple.
func (c Color) ToRGB() (int, int, int) {
	return int(c.Red * 255), int(c.Green * 255), int(c.Blue * 255)
}

// ToHex converts the color to a hexadecimal string.
func (c Color) ToHex(octothorpe bool) string {
	r, g, b := c.ToRGB()
	hex := fmt.Sprintf("%02X%02X%02X", r, g, b)
	if octothorpe {
		return "#" + hex
	}
	return hex
}

// FromHex creates a Color instance from a hexadecimal string.
func FromHex(hex string) (Color, error) {

	hex = trimOctothorpe(hex)

	if len(hex) != 6 && len(hex) != 8 {
		return Color{}, fmt.Errorf("hex color must be 6 or 8 characters long: %s", hex)
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid red component: %w", err)
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid green component: %w", err)
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid blue component: %w", err)
	}

	a := int64(255)
	if len(hex) == 8 {
		a, err = strconv.ParseInt(hex[6:8], 16, 64)
		if err != nil {
			return Color{}, fmt.Errorf("invalid alpha component: %w", err)
		}
	}

	return Color{
		Red:   float64(r) / 255.0,
		Green: float64(g) / 255.0,
		Blue:  float64(b) / 255.0,
		Alpha: float64(a) / 255.0,
	}, nil
}

// ToDict serializes the color to a dictionary-like map.
func (c Color) ToDict() map[string]string {
	return map[string]string{
		"hex": c.ToHex(true),
	}
}

// FromDict deserializes a Color from a dictionary-like map.
func FromDict(data map[string]string) (Color, error) {
	hex, exists := data["hex"]
	if !exists {
		return Color{}, fmt.Errorf("missing 'hex' key in data")
	}
	return FromHex(hex)
}

// Helper function to trim a leading '#' from a hex string.
func trimOctothorpe(hex string) string {
	if len(hex) > 0 && hex[0] == '#' {
		return hex[1:]
	}
	return hex
}

var (
	rnd  *rand.Rand
	once sync.Once
)

func getRand() *rand.Rand {
	once.Do(func() {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return rnd
}

// RandomColor returns a random opaque Color.
func RandomColor() Color {
	r := getRand()
	return Color{
		Red:   r.Float64(),
		Green: r.Float64(),
		Blue:  r.Float64(),
		Alpha: 1.0,
	}
}

func clampFloatToUint8(f float64) uint8 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 255
	}
	return uint8(f * 255)
}

func (c *Color) Hex() string {
	if c == nil {
		return "#ff0000" // default 100% red
	}
	r := clampFloatToUint8(c.Red)
	g := clampFloatToUint8(c.Green)
	b := clampFloatToUint8(c.Blue)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func ColorsSimilar(c1, c2 Color, tol float64) bool {
	dAlpha := c1.Alpha - c2.Alpha
	dRed := c1.Red - c2.Red
	dGreen := c1.Green - c2.Green
	dBlue := c1.Blue - c2.Blue

	distSq := dAlpha*dAlpha + dRed*dRed + dGreen*dGreen + dBlue*dBlue
	return distSq <= tol*tol
}
