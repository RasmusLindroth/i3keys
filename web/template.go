package web

import _ "embed"

//go:embed res/index.gohtml
var indexTmplStr string

//go:embed res/index.css
var indexTmplCSS string

//go:embed res/index.js
var indexTmplJS string
