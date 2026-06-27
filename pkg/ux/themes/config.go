package themes

type ThemeConfig struct {
	SearchOrder []string          // search order for resources for this theme (this lets us have use shared paths for certain objects )
	Colors      map[string]string // maps color name to hex value
	Fonts       map[string]struct {
		Font string
		Size float64
	}
	Images    string
	FontIcons string
}
