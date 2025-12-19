package pagination

import (
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}

	if limit > MaxLimit {
		return MaxLimit
	}

	return limit
}

func EncodeToken(id uint) string {
	if id == 0 {
		return ""
	}

	idStr := strconv.FormatUint(uint64(id), 10)
	idBytes := []byte(idStr)
	paddingSize := 6
	padding := make([]byte, paddingSize)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range padding {
		padding[i] = byte(rng.Intn(256))
	}

	combined := append(padding, idBytes...)

	return base64.URLEncoding.EncodeToString(combined)
}

func DecodeToken(token string) (uint, error) {
	if token == "" {
		return 0, nil
	}

	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}

	paddingSize := 6
	if len(data) <= paddingSize {
		return 0, base64.CorruptInputError(0)
	}

	idBytes := data[paddingSize:]
	id, err := strconv.ParseUint(string(idBytes), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func ParseLimit(limitStr string) int {
	if limitStr == "" {
		return DefaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return DefaultLimit
	}

	return NormalizeLimit(limit)
}

func FormatTokenForJSON(token string) interface{} {
	if token == "" {
		return nil
	}

	return token
}

func BuildPaginatedResponse(items interface{}, token string) map[string]interface{} {
	return map[string]interface{}{
		"items":     items,
		"nextToken": FormatTokenForJSON(token),
	}
}
