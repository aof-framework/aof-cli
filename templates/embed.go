package templates

import "embed"

// FS contains the pinned, local AOF bootstrap assets used by aof init.
//
//go:embed static/** dynamic/**
var FS embed.FS
