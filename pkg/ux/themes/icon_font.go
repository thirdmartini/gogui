package themes

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

type Icon interface {
	Draw(canvas canvas.Canvas, x, y int, color color.Color)
}

type IconProvider interface {
	GetIcon(name string) Icon
}

type IconNullop struct {
}

func (in *IconNullop) Draw(canvas canvas.Canvas, x, y int, color color.Color) {
}

type IconFontProvider struct {
	font *fonts.Font
	code map[string]rune
}

func (ifp *IconFontProvider) GetIcon(name string) Icon {
	r, ok := ifp.code[name]
	if !ok {
		panic(fmt.Sprintf("icon %s not found", name))
		return &IconNullop{}
	}

	fw, fh := ifp.font.Measure(string(r))

	return &IconFont{
		code: r,
		font: ifp.font,
		xoff: -fw / 2,
		yoff: fh / 2,
	}
}

func NewIconFontProvider(fontPath string, points float64) (*IconFontProvider, error) {
	font, err := fonts.Load(fontPath, points)
	if err != nil {
		panic(err)
		return nil, err
	}

	codepointsPath := strings.TrimSuffix(fontPath, path.Ext(fontPath))

	codes, err := parseIconFile(codepointsPath + ".codepoints")
	if err != nil {
		panic(err)
		return nil, err
	}

	return &IconFontProvider{
		font: font,
		code: codes,
	}, nil
}

type IconFont struct {
	code       rune
	font       *fonts.Font
	xoff, yoff int
}

func (ifi *IconFont) Draw(canvas canvas.Canvas, x, y int, color color.Color) {
	canvas.DrawText(x+ifi.xoff, y+ifi.yoff, string(ifi.code), ifi.font, color)
}

func parseIconFile(filePath string) (map[string]rune, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	iconMap := make(map[string]rune)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { // skip empty lines and comments if any
			continue
		}

		// Split on whitespace (handles multiple spaces)
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("Warning: line %d has insufficient fields: %s\n", lineNum, line)
			continue
		}

		name := parts[0]
		hexStr := parts[1]

		// Convert hex string to int (base 16)
		val, err := strconv.ParseUint(hexStr, 16, 64)
		if err != nil {
			fmt.Printf("Warning: line %d invalid hex '%s': %v\n", lineNum, hexStr, err)
			continue
		}

		// Use int (assuming values fit; use int64 if needed for larger values)
		iconMap[name] = rune(val)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return iconMap, nil
}
