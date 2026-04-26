// Package embed provides embedded dictionary data for jiebago.
// This allows jiebago to work out-of-the-box without external dictionary files.
package embed

import (
	_ "embed"
)

//go:embed dict.txt
var DictData []byte
