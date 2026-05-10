package tui

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
)

type Theme struct {
	TitleBG        color.Color
	TitleText      color.Color
	KeyInUse       color.Color
	KeyNormal      color.Color
	ListTitle      color.Color
	NormalTitle    color.Color
	NormalDesc     color.Color
	SelectedTitle  color.Color
	SelectedDesc   color.Color
	SelectedBorder color.Color
	DimmedTitle    color.Color
	DimmedDesc     color.Color
}

var LightTheme = Theme{
	TitleBG:        lipgloss.Black,
	TitleText:      lipgloss.White,
	KeyInUse:       lipgloss.Color("#CC3333"),
	KeyNormal:      lipgloss.Color("#222222"),
	ListTitle:      lipgloss.Color("#111111"),
	NormalTitle:    lipgloss.Color("#111111"),
	NormalDesc:     lipgloss.Color("#444444"),
	SelectedTitle:  lipgloss.Color("#000000"),
	SelectedDesc:   lipgloss.Color("#333333"),
	SelectedBorder: lipgloss.Color("#70b5a7"),
	DimmedTitle:    lipgloss.Color("#666666"),
	DimmedDesc:     lipgloss.Color("#888888"),
}

var DarkTheme = Theme{
	TitleBG:        lipgloss.White,
	TitleText:      lipgloss.Black,
	KeyInUse:       lipgloss.Color("#FF6B6B"),
	KeyNormal:      lipgloss.Color("#DDDDDD"),
	ListTitle:      lipgloss.Color("#FFFFFF"),
	NormalTitle:    lipgloss.Color("#EEEEEE"),
	NormalDesc:     lipgloss.Color("#AAAAAA"),
	SelectedTitle:  lipgloss.Color("#FFFFFF"),
	SelectedDesc:   lipgloss.Color("#CCCCCC"),
	SelectedBorder: lipgloss.Color("#a2e8da"),
	DimmedTitle:    lipgloss.Color("#777777"),
	DimmedDesc:     lipgloss.Color("#555555"),
}

func ActiveTheme() Theme {
	if lipgloss.HasDarkBackground(os.Stdin, os.Stdout) {
		return DarkTheme
	}
	return LightTheme
}
