package assets

import "embed"

// FileSystem holds Management UI template files and its assets
//
//go:embed *.html
var FileSystem embed.FS
