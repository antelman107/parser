package top

import "github.com/migotom/heavykeeper"

// HkConfigDefault is standard config.
var HkConfigDefault = HkConfig{
	Workers: 1, // results seem to be random if this is > 1
	Width:   2048,
	Depth:   5,
	Decay:   0.9,
}

// HkConfig defines heavy keeper settings besides K.
type HkConfig struct {
	Workers int
	Width   uint32
	Depth   uint32
	Decay   float64
}

// GetHK returns initialized TopK by top length (k) and cfg struct.
func GetHK(k uint32, cfg HkConfig) *heavykeeper.TopK {
	return heavykeeper.New(cfg.Workers, k, cfg.Width, cfg.Depth, cfg.Decay)
}
