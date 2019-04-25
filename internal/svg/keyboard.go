package svg

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"text/template"

	"github.com/RasmusLindroth/i3keys/internal/keyboard"
)

var tmpl = `<svg width="1500" height="450" xmlns="http://www.w3.org/2000/svg">
{{- $currY := 10.0 -}}
{{- range $y, $yItem := . -}}
{{- $currX := 10.0 -}}
{{- $currK := 0 -}}
	{{- range $x, $xItem := $yItem -}}
		{{- if $xItem.Visible -}}
			{{ if $xItem.Enter }}
			<g>
				<path stroke="black" fill="{{fill $xItem}}" d="{{enter $currX $currY $xItem}}" />
				{{text $currX $currY $xItem}}
			</g>
			{{- else -}}
			<g>
				<rect x="{{$currX}}" y="{{$currY}}" width="{{widthInner $xItem}}" height="{{height $xItem}}" stroke="black" fill="{{fill $xItem}}" />
				{{text $currX $currY $xItem}}
			</g>
			{{- end -}}
			{{- $currK = incOne $currK -}}
		{{- end -}}
		{{- $currX = incX $currX $xItem -}}
	{{- end -}}
{{- $currY = incY $currY -}}
{{- end -}}
</svg>
`

type keyType struct {
	X        float64
	Y        float64
	Visible  bool
	Enter    bool
	Text     string
	Modifier bool
	InUse    bool
}

var (
	single      = keyType{1, 1, true, false, "", false, false}    //one key
	emptySingle = keyType{1, 1, false, false, "", false, false}   //empty one key
	emptySmall  = keyType{0.5, 1, false, false, "", false, false} //gap between f-keys
	double      = keyType{2, 1, true, false, "", false, false}    //backspace
	onehalf     = keyType{1.5, 1, true, false, "", false, false}  //tab
	doubleY     = keyType{1, 2, true, false, "", false, false}    //numpad +, enter
	semidouble  = keyType{1.75, 1, true, false, "", false, false} //caps
	modifier    = keyType{1.25, 1, true, false, "", false, false} //ctrl
	semilarge   = keyType{2.25, 1, true, false, "", false, false}
	large       = keyType{2.75, 1, true, false, "", false, false}  //right shift
	space       = keyType{6.25, 1, true, false, "", false, false}  //space
	enterISO    = keyType{1.5, 1, true, true, "", false, false}    //enter
	enterDown   = keyType{1.25, 1, false, false, "", false, false} //enter
)

var ansi = [][]keyType{
	[]keyType{single, emptySingle, single, single, single, single, emptySmall, single, single, single, single, emptySmall, single, single, single, single, emptySmall, single, single, single},
	[]keyType{single, single, single, single, single, single, single, single, single, single, single, single, single, double, emptySmall, single, single, single, emptySmall, single, single, single, single},
	[]keyType{onehalf, single, single, single, single, single, single, single, single, single, single, single, single, onehalf, emptySmall, single, single, single, emptySmall, single, single, single, doubleY},
	[]keyType{semidouble, single, single, single, single, single, single, single, single, single, single, single, semilarge, emptySmall, emptySingle, emptySingle, emptySingle, emptySmall, single, single, single},
	[]keyType{semilarge, single, single, single, single, single, single, single, single, single, single, large, emptySmall, emptySingle, single, emptySingle, emptySmall, single, single, single, doubleY},
	[]keyType{modifier, modifier, modifier, space, modifier, modifier, modifier, modifier, emptySmall, single, single, single, emptySmall, double, single},
}

var iso = [][]keyType{
	[]keyType{single, emptySingle, single, single, single, single, emptySmall, single, single, single, single, emptySmall, single, single, single, single, emptySmall, single, single, single},
	[]keyType{single, single, single, single, single, single, single, single, single, single, single, single, single, double, emptySmall, single, single, single, emptySmall, single, single, single, single},
	[]keyType{onehalf, single, single, single, single, single, single, single, single, single, single, single, single, enterISO, emptySmall, single, single, single, emptySmall, single, single, single, doubleY},
	[]keyType{semidouble, single, single, single, single, single, single, single, single, single, single, single, single, enterDown, emptySmall, emptySingle, emptySingle, emptySingle, emptySmall, single, single, single},
	[]keyType{modifier, single, single, single, single, single, single, single, single, single, single, single, large, emptySmall, emptySingle, single, emptySingle, emptySmall, single, single, single, doubleY},
	[]keyType{modifier, modifier, modifier, space, modifier, modifier, modifier, modifier, emptySmall, single, single, single, emptySmall, double, single},
}

//Generate creates an SVG image of the keyboard
func Generate(layout string, kb keyboard.Keyboard) []byte {
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"incX": func(currX float64, key keyType) (n float64) {
			return currX + (key.X * 64)
		},
		"incY": func(currY float64) (n float64) {
			return currY + 70
		},
		"incOne": func(currN int) int {
			return currN + 1
		},
		"width": func(key keyType) (n float64) {
			return key.X * 64
		},
		"height": func(key keyType) (n float64) {
			if key.Y > 1.0 {
				return (key.Y * 64) + 6
			}
			return key.Y * 64
		},
		"widthInner": func(key keyType) (n float64) {
			return (key.X * 64) - 6
		},
		"enter": func(currX float64, currY float64, key keyType) string {
			topXEnd := currX + (1.5 * 64) - 6
			bottomY := currY + (2.0 * 64) + 6
			bottomXStart := currX + (0.25 * 64)
			return fmt.Sprintf("M%.2f,%.2f L%.2f,%.2f L%.2f,%.2f L%.2f,%.2f L%.2f,%.2f L%.2f,%.2f L%.2f,%.2f",
				currX, currY, topXEnd, currY, topXEnd, bottomY, bottomXStart, bottomY, bottomXStart, currY+64, currX, currY+64, currX, currY)
		},
		"fill": func(key keyType) string {
			if key.Modifier {
				return "#5cff87"
			}

			if key.InUse {
				return "#ff655c"
			}

			return "none"
		},
		"text": func(currX float64, currY float64, key keyType) string {
			txt := key.Text
			charPerRow := int(key.X * 7)
			rows := float64(len(txt)) / float64(charPerRow)
			startX := currX + 2 + (key.X * 32) - 6
			startY := currY - 3

			if rows <= 1.0 {

				return fmt.Sprintf("<text dominant-baseline=\"middle\" text-anchor=\"middle\" x=\"%.0f\" y=\"%.0f\" font-family=\"Monospace\" font-size=\"12\"  fill=\"black\">%s</text>",
					startX, startY+(12*(3*key.Y)), txt)
			}

			r := fmt.Sprintf("<text dominant-baseline=\"middle\" text-anchor=\"middle\" x=\"%.0f\" y=\"%.0f\" font-family=\"Monospace\" font-size=\"12\"  fill=\"black\">",
				startX, startY+(6*((5*key.Y)-float64(math.Ceil(rows)))))

			j := 0
			for i := 0.0; i < rows; i = i + 1.0 {
				startTxt := j * charPerRow
				endTxt := charPerRow + (j * charPerRow)
				if endTxt > len(txt) {
					endTxt = len(txt)
				}
				txtTmp := txt[startTxt:endTxt]
				r = fmt.Sprintf("%s<tspan x=\"%.0f\" dy=\"%.0f\">%s</tspan>", r, startX, 12.0, txtTmp)
				j++
			}
			return fmt.Sprintf("%s</text>", r)
		},
	}

	templ, err := template.New("").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	usedKb := iso
	if layout == "ANSI" {
		usedKb = ansi
	}

	for i := 0; i < len(usedKb); i++ {
		k := 0
		for j := 0; j < len(usedKb[i]); j++ {
			if usedKb[i][j].Visible == false {
				continue
			}

			currKey := kb.Keys[i][k]
			usedKb[i][j].Text = currKey.Symbol
			usedKb[i][j].Modifier = currKey.Modifier
			usedKb[i][j].InUse = currKey.InUse

			k++
		}
	}

	var data bytes.Buffer
	err = templ.Execute(&data, usedKb)
	if err != nil {
		log.Fatalln(err)
	}
	return data.Bytes()
}
