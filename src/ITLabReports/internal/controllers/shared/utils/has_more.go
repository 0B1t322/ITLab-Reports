package utils

func HasMore(
	totalResult,
	limit,
	offset int64,
) bool {
	return limit != 0 && totalResult-offset-limit > 0
}
