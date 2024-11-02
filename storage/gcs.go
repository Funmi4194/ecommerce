package storage

import (
	"github.com/funmi4194/ecommerce/primer"
	"github.com/opensaucerer/bifrost"
)

// NewGCSRainbowBridge creates a new bifrost rainbow bridge to  cloud storage
func NewGCSRainbowBridge(bucket string) (bifrost.RainbowBridge, error) {
	// database.RainbowBridge
	return bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket: bucket,
		// DefaultTimeout:  10,
		Provider:        bifrost.GoogleCloudStorage,
		EnableDebug:     true,
		PublicRead:      false,
		CredentialsFile: primer.ENV.GoogleApplicationCredentials,
	})
}
