package model

type CreatePacksRequest struct {
	Packs []int64 `json:"packs"`
}

type CreatePacksResponse struct {
	VersionHash string `json:"version_hash"`
}

type PacksListResponse struct {
	Packs       []int64 `json:"packs"`
	VersionHash string  `json:"version_hash"`
}
