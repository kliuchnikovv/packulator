package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePacksRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		request  CreatePacksRequest
		expected string
	}{
		{
			name: "empty packs",
			request: CreatePacksRequest{
				Packs: []int64{},
			},
			expected: `{"packs":[]}`,
		},
		{
			name: "single pack",
			request: CreatePacksRequest{
				Packs: []int64{250},
			},
			expected: `{"packs":[250]}`,
		},
		{
			name: "multiple packs",
			request: CreatePacksRequest{
				Packs: []int64{250, 500, 1000, 2000, 5000},
			},
			expected: `{"packs":[250,500,1000,2000,5000]}`,
		},
		{
			name: "packs with zero",
			request: CreatePacksRequest{
				Packs: []int64{0, 250, 500},
			},
			expected: `{"packs":[0,250,500]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.request)
			require.NoError(t, err)
			
			assert.JSONEq(t, tt.expected, string(jsonData))
		})
	}
}

func TestCreatePacksRequest_JSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected CreatePacksRequest
	}{
		{
			name:     "empty packs",
			jsonData: `{"packs":[]}`,
			expected: CreatePacksRequest{
				Packs: []int64{},
			},
		},
		{
			name:     "single pack",
			jsonData: `{"packs":[250]}`,
			expected: CreatePacksRequest{
				Packs: []int64{250},
			},
		},
		{
			name:     "multiple packs",
			jsonData: `{"packs":[250,500,1000,2000,5000]}`,
			expected: CreatePacksRequest{
				Packs: []int64{250, 500, 1000, 2000, 5000},
			},
		},
		{
			name:     "packs with zero",
			jsonData: `{"packs":[0,250,500]}`,
			expected: CreatePacksRequest{
				Packs: []int64{0, 250, 500},
			},
		},
		{
			name:     "null packs field",
			jsonData: `{"packs":null}`,
			expected: CreatePacksRequest{
				Packs: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request CreatePacksRequest
			err := json.Unmarshal([]byte(tt.jsonData), &request)
			require.NoError(t, err)
			
			assert.Equal(t, tt.expected, request)
		})
	}
}

func TestCreatePacksRequest_InvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON syntax",
			jsonData: `{"packs":[250,500,}`,
		},
		{
			name:     "wrong type for packs",
			jsonData: `{"packs":"invalid"}`,
		},
		{
			name:     "wrong type for pack values",
			jsonData: `{"packs":["250","500"]}`,
		},
		{
			name:     "missing packs field",
			jsonData: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request CreatePacksRequest
			err := json.Unmarshal([]byte(tt.jsonData), &request)
			
			if tt.name == "missing packs field" {
				// Missing field should not error, but result in nil slice
				assert.NoError(t, err)
				assert.Nil(t, request.Packs)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCreatePacksResponse_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		response CreatePacksResponse
		expected string
	}{
		{
			name: "valid version hash",
			response: CreatePacksResponse{
				VersionHash: "abc123def456",
			},
			expected: `{"version_hash":"abc123def456"}`,
		},
		{
			name: "empty version hash",
			response: CreatePacksResponse{
				VersionHash: "",
			},
			expected: `{"version_hash":""}`,
		},
		{
			name: "long version hash",
			response: CreatePacksResponse{
				VersionHash: "1234567890abcdef1234567890abcdef12345678",
			},
			expected: `{"version_hash":"1234567890abcdef1234567890abcdef12345678"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.response)
			require.NoError(t, err)
			
			assert.JSONEq(t, tt.expected, string(jsonData))
		})
	}
}

func TestCreatePacksResponse_JSONUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected CreatePacksResponse
	}{
		{
			name:     "valid version hash",
			jsonData: `{"version_hash":"abc123def456"}`,
			expected: CreatePacksResponse{
				VersionHash: "abc123def456",
			},
		},
		{
			name:     "empty version hash",
			jsonData: `{"version_hash":""}`,
			expected: CreatePacksResponse{
				VersionHash: "",
			},
		},
		{
			name:     "missing version_hash field",
			jsonData: `{}`,
			expected: CreatePacksResponse{
				VersionHash: "",
			},
		},
		{
			name:     "null version_hash field",
			jsonData: `{"version_hash":null}`,
			expected: CreatePacksResponse{
				VersionHash: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response CreatePacksResponse
			err := json.Unmarshal([]byte(tt.jsonData), &response)
			require.NoError(t, err)
			
			assert.Equal(t, tt.expected, response)
		})
	}
}

func TestCreatePacksResponse_InvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name:     "invalid JSON syntax",
			jsonData: `{"version_hash":"abc123",}`,
		},
		{
			name:     "wrong type for version_hash",
			jsonData: `{"version_hash":123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response CreatePacksResponse
			err := json.Unmarshal([]byte(tt.jsonData), &response)
			assert.Error(t, err)
		})
	}
}

func TestCreatePacksRequest_ZeroValues(t *testing.T) {
	var request CreatePacksRequest
	
	assert.Nil(t, request.Packs)
}

func TestCreatePacksResponse_ZeroValues(t *testing.T) {
	var response CreatePacksResponse
	
	assert.Equal(t, "", response.VersionHash)
}

func TestCreatePacksRequest_RoundTrip(t *testing.T) {
	original := CreatePacksRequest{
		Packs: []int64{250, 500, 1000, 2000, 5000},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal from JSON
	var unmarshaled CreatePacksRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Should be equal to original
	assert.Equal(t, original, unmarshaled)
}

func TestCreatePacksResponse_RoundTrip(t *testing.T) {
	original := CreatePacksResponse{
		VersionHash: "abc123def456789",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal from JSON
	var unmarshaled CreatePacksResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	// Should be equal to original
	assert.Equal(t, original, unmarshaled)
}