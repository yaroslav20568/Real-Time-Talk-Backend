package pagination

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type TokenData struct {
	ID        uint
	Timestamp int64
}

func GenerateToken(id uint, timestamp int64) string {
	tokenData := fmt.Sprintf("%d:%d", id, timestamp)
	return base64.StdEncoding.EncodeToString([]byte(tokenData))
}

func ParseToken(token string) (*TokenData, error) {
	if token == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid nextToken: %w", err)
	}

	tokenParts := string(decoded)
	var id uint64
	var timestamp int64

	if n, err := fmt.Sscanf(tokenParts, "%d:%d", &id, &timestamp); err == nil && n == 2 {
		return &TokenData{
			ID:        uint(id),
			Timestamp: timestamp,
		}, nil
	}

	if id, err := strconv.ParseUint(tokenParts, 10, 32); err == nil {
		return &TokenData{
			ID: uint(id),
		}, nil
	}

	return nil, fmt.Errorf("invalid nextToken format")
}

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}
