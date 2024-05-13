package main

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/image/font/inconsolata"
)

const (
	screenWidth  = 777
	screenHeight = 388

	toolbarHeight      = 40
	sliderWidth        = 200
	sliderHeight       = 10
	sliderKnobDiameter = 20
	toggleDiameter     = 20
	pointDiameter      = 10
)

var (
	lavender = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Lavender().Hex,
		Dark:  catppuccin.Mocha.Lavender().Hex,
	}
	peach = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Peach().Hex,
		Dark:  catppuccin.Mocha.Peach().Hex,
	}
	blue = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Blue().Hex,
		Dark:  catppuccin.Mocha.Blue().Hex,
	}
	pink = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Pink().Hex,
		Dark:  catppuccin.Mocha.Pink().Hex,
	}
	surface0 = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Surface0().Hex,
		Dark:  catppuccin.Mocha.Surface0().Hex,
	}
	surface1 = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Surface1().Hex,
		Dark:  catppuccin.Mocha.Surface1().Hex,
	}
	surface2 = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Surface2().Hex,
		Dark:  catppuccin.Mocha.Surface2().Hex,
	}
	green = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Green().Hex,
		Dark:  catppuccin.Mocha.Green().Hex,
	}
	red = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Red().Hex,
		Dark:  catppuccin.Mocha.Red().Hex,
	}
	yellow = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Yellow().Hex,
		Dark:  catppuccin.Mocha.Yellow().Hex,
	}
	mauve = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Mauve().Hex,
		Dark:  catppuccin.Mocha.Mauve().Hex,
	}
	rosewater = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Rosewater().Hex,
		Dark:  catppuccin.Mocha.Rosewater().Hex,
	}
	sky = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Sky().Hex,
		Dark:  catppuccin.Mocha.Sky().Hex,
	}
	teal = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Teal().Hex,
		Dark:  catppuccin.Mocha.Teal().Hex,
	}
	flamingo = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Flamingo().Hex,
		Dark:  catppuccin.Mocha.Flamingo().Hex,
	}
	maroon = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Maroon().Hex,
		Dark:  catppuccin.Mocha.Maroon().Hex,
	}
	sapphire = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Sapphire().Hex,
		Dark:  catppuccin.Mocha.Sapphire().Hex,
	}
	crust = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Crust().Hex,
		Dark:  catppuccin.Mocha.Crust().Hex,
	}
	mantle = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Mantle().Hex,
		Dark:  catppuccin.Mocha.Mantle().Hex,
	}
	overlay1 = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Overlay1().Hex,
		Dark:  catppuccin.Mocha.Overlay1().Hex,
	}
	overlay0 = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Overlay0().Hex,
		Dark:  catppuccin.Mocha.Overlay0().Hex,
	}
	overlay2 = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Overlay2().Hex,
		Dark:  catppuccin.Mocha.Overlay2().Hex,
	}
	base = lipgloss.AdaptiveColor{
		Light: catppuccin.Latte.Base().Hex,
		Dark:  catppuccin.Mocha.Base().Hex,
	}
	colors = []lipgloss.AdaptiveColor{
		rosewater,
		flamingo,
		pink,
		mauve,
		red,
		maroon,
		peach,
		yellow,
		green,
		teal,
		sky,
		sapphire,
		blue,
		lavender,
	}

	textFont = inconsolata.Regular8x16

	backgroundColor   = green
	toolbarColor      = mauve
	sliderBgColor     = surface1
	sliderKnobColor   = sky
	toggleOnColor     = sky
	toggleOffColor    = maroon
	textColor         = crust
	pointColor        = red
	outlineColor      = sliderBgColor
	curveColor        = maroon
	naturalCurveColor = overlay2

	padding = sliderKnobDiameter
)
