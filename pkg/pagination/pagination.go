package pagination

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

func NormalizePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func CalculateOffset(page int, limit int) int {
	normalizedPage := NormalizePage(page)
	return (normalizedPage - 1) * limit
}

func CalculateTotalPages(total int64, limit int) int {
	if total <= 0 {
		return 1
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		return 1
	}
	return totalPages
}
