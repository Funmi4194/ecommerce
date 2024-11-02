package storage

import (
	"github.com/funmi4194/bifrost"
	"github.com/funmi4194/ecommerce/primer"
)

// NewGCSRainbowBridge creates a new bifrost rainbow bridge to  cloud storage
func NewGCSRainbowBridge(bucket string) (bifrost.RainbowBridge, error) {
	return bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:   bucket,
		DefaultTimeout:  10,
		Provider:        bifrost.GoogleCloudStorage,
		EnableDebug:     true,
		PublicRead:      false,
		CredentialsFile: primer.ENV.GoogleApplicationCredentials,
	})
}
