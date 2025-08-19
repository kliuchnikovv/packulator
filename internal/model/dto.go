// Package model contains data transfer objects (DTOs) for API communication.
package model

// CreatePacksRequest represents the payload for creating a new pack configuration.
// It contains an array of pack sizes that will be available for packaging calculations.
type CreatePacksRequest struct {
	Packs []int64 `json:"packs"` // Array of available pack sizes
}

// CreatePacksResponse represents the response after creating a pack configuration.
// It returns the version hash that can be used to reference this pack configuration.
type CreatePacksResponse struct {
	VersionHash string `json:"version_hash"` // Unique hash identifying the pack configuration
}
