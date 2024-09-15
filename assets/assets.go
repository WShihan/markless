package assets

import "embed"

var (
	//go:embed html/*.html
	HTML embed.FS
	//go:embed static/js/*
	JS embed.FS
	//go:embed static/css/*
	CSS embed.FS
	// //go:embed static/img/*
	IMG embed.FS
	// //go:embed static/img/*.ico
	ICO embed.FS
)
