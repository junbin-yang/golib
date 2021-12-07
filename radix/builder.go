package radix

import (
	"github.com/gbrlsnchs/color"
	"strings"
)

const (
	colorRed = iota
	colorGreen
	colorMagenta
	colorBold
)

type builder struct {
	*strings.Builder
	colors [4]color.Color
	debug  bool
}
