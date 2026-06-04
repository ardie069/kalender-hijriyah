package kernels

import "embed"

//go:embed *.bsp *.tls *.tpc
var FS embed.FS
