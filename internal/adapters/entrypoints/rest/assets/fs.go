package assets

import "embed"

// TemplateFileSystem holds Management UI template files
//
//go:embed *.html
var TemplateFileSystem embed.FS

// FileSystem holds Management UI template assets (images, styles, scripts)
//
//go:embed static favicon.ico
var FileSystem embed.FS
