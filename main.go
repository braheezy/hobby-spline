package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/braheezy/hobby-spline/pkg/bezier"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hobby's algorithm for aesthetic Bézier splines")
	game := &Game{
		points: []bezier.Point{
			{X: 356, Y: 229},
			{X: 523, Y: 287},
			{X: 505, Y: 72},
			{X: 109, Y: 224},
			{X: 108, Y: 92},
			{X: 232, Y: 307},
		},
		omega:       0.75,
		showComb:    true,
		showNatural: true,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	points         []bezier.Point
	splinePoints   []bezier.Point
	omega          float64
	showComb       bool
	showNatural    bool
	sliderDragging bool
	draggingPoint  *bezier.Point
	dragOffsetX    float32
	dragOffsetY    float32
}

func (g *Game) Update() error {

	// Handle slider for omega
	mx, my := ebiten.CursorPosition()
	// Calculate UI control positions
	sliderX := (screenWidth - sliderWidth) / 2
	sliderY := screenHeight - toolbarHeight/2 - sliderHeight/2
	toggleX := sliderX + sliderWidth + 60 // 60 pixels away from the slider end
	toggleY := screenHeight - toolbarHeight/2 - toggleDiameter/2

	// Check mouse interactions with slider knob
	if g.sliderDragging {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			g.sliderDragging = false
		} else {
			g.omega = float64(mx-sliderX) / float64(sliderWidth)
			if g.omega < 0 {
				g.omega = 0
			} else if g.omega > 1 {
				g.omega = 1
			}
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		knobX := sliderX + int(g.omega*float64(sliderWidth))
		// Check if mouse is within slider knob
		if mx >= knobX-sliderKnobDiameter/2 && mx <= knobX+sliderKnobDiameter/2 && my >= sliderY && my <= sliderY+sliderHeight {
			g.sliderDragging = true
		}
		// Check if mouse is within toggle
		if mx >= toggleX-toggleDiameter/2 && mx <= toggleX+toggleDiameter/2 && my >= toggleY && my <= toggleY+toggleDiameter {
			g.showComb = !g.showComb
		}
	}

	// Handle point dragging logic for bezier curves
	if g.draggingPoint != nil {
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			g.draggingPoint = nil
		} else {
			x, y := ebiten.CursorPosition()
			g.draggingPoint.X = float64(x) - float64(g.dragOffsetX)
			g.draggingPoint.Y = float64(y) - float64(g.dragOffsetY)
		}
	} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for i := range g.points {
			px, py := g.points[i].X, g.points[i].Y
			if (x-int(px))*(x-int(px))+(y-int(py))*(y-int(py)) <= 25*25 { // 25 pixels radius for easier clicking
				g.draggingPoint = &g.points[i]
				g.dragOffsetX = float32(x) - float32(px)
				g.dragOffsetY = float32(y) - float32(py)
				break
			}
		}
	}

	g.splinePoints, _ = bezier.CreateHobbySpline(g.points, g.omega)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Set background
	screen.Fill(backgroundColor)
	if g.showNatural {
		// Calculate natural spline
		naturalPoints, _ := bezier.NaturalCubicSpline(g.points)
		if len(naturalPoints) != 0 {
			for i := 0; i <= (len(naturalPoints)-2)/3; i++ {
				pts := naturalPoints[i*3 : i*3+4]
				b := &bezier.Bezier{Points: pts}
				strokeCurve(screen, b)
			}
		}
	}
	// Draw bezier curves
	if len(g.splinePoints) != 0 {
		for i := 0; i <= (len(g.splinePoints)-2)/3; i++ {
			pts := g.splinePoints[i*3 : i*3+4]
			curve, err := bezier.NewBezier(false, pts...)
			if err != nil {
				log.Fatal(err)
			}
			if g.showComb {
				drawComb(screen, curve)
			}
			strokeCurve(screen, curve)
		}
	}

	// Draw points that user can grab
	for _, pt := range g.points {
		vector.DrawFilledCircle(screen, float32(pt.X), float32(pt.Y), pointDiameter, pointColor, true)
	}

	// UI elements
	// Draw the toolbar
	vector.DrawFilledRect(screen, 0, float32(screenHeight-toolbarHeight), float32(screenWidth), float32(toolbarHeight), toolbarColor, true)

	// Calculate positions
	sliderX := (screenWidth - sliderWidth) / 2
	toggleX := sliderX + sliderWidth + 60
	sliderY := screenHeight - toolbarHeight/2 - sliderHeight/2
	toggleY := screenHeight - toolbarHeight/2

	// Draw slider background
	vector.DrawFilledRect(screen, float32(sliderX), float32(sliderY), float32(sliderWidth), float32(sliderHeight), sliderBgColor, true)

	// Draw slider knob
	knobX := sliderX + int(g.omega*float64(sliderWidth))
	vector.DrawFilledCircle(screen, float32(knobX), float32(screenHeight-toolbarHeight/2), float32(sliderKnobDiameter/2), sliderKnobColor, true)
	vector.StrokeCircle(screen, float32(knobX), float32(screenHeight-toolbarHeight/2), float32(sliderKnobDiameter/2), 1, outlineColor, true)

	// Draw toggle
	toggleColor := toggleOnColor
	if !g.showComb {
		toggleColor = toggleOffColor
	}
	vector.DrawFilledCircle(screen, float32(toggleX), float32(toggleY), float32(sliderHeight), toggleColor, true)
	vector.StrokeCircle(screen, float32(toggleX), float32(toggleY), float32(sliderHeight), 1, outlineColor, true)

	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(textColor)
	// Draw omega value
	omegaText := fmt.Sprintf("w = %.2f", g.omega)
	omegaTextWidth := textWidth(omegaText, textFont)
	textOp.GeoM.Translate(float64(sliderX-omegaTextWidth-padding), float64(sliderY-3))
	text.Draw(screen, omegaText, text.NewGoXFace(textFont), textOp)

	// Draw "Show Comb" label next to toggle
	textOp = &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(textColor)
	textOp.GeoM.Translate(float64(toggleX+toggleDiameter), float64(toggleY-sliderHeight/2-3))
	text.Draw(screen, "Show Comb", text.NewGoXFace(textFont), textOp)

}

func textWidth(s string, face font.Face) int {
	bounds, _ := font.BoundString(face, s)
	return (bounds.Max.X - bounds.Min.X).Ceil()
}

var cachedColors map[int][]color.Color = make(map[int][]color.Color)

func getCombColors(teeth int) []color.Color {
	if colors, found := cachedColors[teeth]; found {
		return colors
	}

	// Calculate new colors and cache them
	newColors := make([]color.Color, teeth)
	for i := 0; i < teeth; i++ {
		newColors[i] = getRainbowColor(float64(i), teeth)
	}
	cachedColors[teeth] = newColors
	return newColors
}

// Function to interpolate a rainbow color based on position
func getRainbowColor(position float64, total int) color.Color {
	if len(colors) == 0 {
		return color.RGBA{0, 0, 0, 255} // return black if no colors are available
	}
	if total <= 1 {
		// If total is 1 or less, return the first color to avoid division by zero
		return colors[0]
	}

	// Ensure position is within [0, total-1]
	if position < 0 {
		position = 0
	} else if position >= float64(total) {
		position = float64(total - 1)
	}

	// Scale position to range within [0, len(colors) - 1]
	scale := (position / float64(total-1)) * float64(len(colors)-1)
	index := int(scale)
	fraction := scale - float64(index)

	if index >= len(colors)-1 {
		return colors[len(colors)-1]
	}

	// Interpolate between colors[index] and colors[index+1]
	c1, c2 := colors[index], colors[index+1]
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	r := interpolateColorComponent(r1, r2, fraction)
	g := interpolateColorComponent(g1, g2, fraction)
	b := interpolateColorComponent(b1, b2, fraction)

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func interpolateColorComponent(c1, c2 uint32, fraction float64) uint8 {
	return uint8((float64(c1>>8)*(1-fraction) + float64(c2>>8)*fraction))
}

func strokeCurve(dst *ebiten.Image, curve *bezier.Bezier) {
	path := &vector.Path{}
	path.MoveTo(float32(curve.Points[0].X), float32(curve.Points[0].Y))
	if len(curve.Points) == 3 {
		path.QuadTo(float32(curve.Points[1].X), float32(curve.Points[1].Y), float32(curve.Points[2].X), float32(curve.Points[2].Y))
	} else if len(curve.Points) == 4 {
		path.CubicTo(float32(curve.Points[1].X), float32(curve.Points[1].Y), float32(curve.Points[2].X), float32(curve.Points[2].Y), float32(curve.Points[3].X), float32(curve.Points[3].Y))
	}
	strokeOp := &vector.StrokeOptions{Width: 5}
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, strokeOp)
	drawVerticesForUtil(dst, vs, is, curveColor, true)
	path.Close()
}

const PIXELS_PER_COMB_TOOTH = 8

func drawComb(dst *ebiten.Image, curve *bezier.Bezier) {
	length := curve.Length()
	teeth := math.Floor(length / PIXELS_PER_COMB_TOOTH)
	colors := getCombColors(int(teeth))
	step := 1 / teeth
	for i := 0.0; i < teeth; i++ {
		t := i * step
		p := curve.Get(t)
		n := curve.Normal(t)
		kr := curve.Curvature(t)
		p2 := bezier.Point{X: p.X + n.X*kr.K*-1500, Y: p.Y + n.Y*kr.K*-1500}
		combColor := colors[int(i)]
		// combColor := red
		vector.StrokeLine(dst, float32(p.X), float32(p.Y), float32(p2.X), float32(p2.Y), 2, combColor, true)
	}
}

func drawVerticesForUtil(dst *ebiten.Image, vs []ebiten.Vertex, is []uint16, clr color.Color, antialias bool) {
	r, g, b, a := clr.RGBA()
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	op.Filter = ebiten.FilterNearest
	op.FillRule = ebiten.NonZero
	op.AntiAlias = antialias
	dst.DrawTriangles(vs, is, whiteSubImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}
