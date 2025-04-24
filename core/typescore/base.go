package typescore

// объединяет общие параметры запроса.
type ListDbOptions struct {
	Filtering  interface{}
	LikeFields map[string]string
	Offset     *uint64
	Limit      *uint64
}

type InsertOptions struct {
	Prefix         string // Префикс для запроса (например, WITH)
	Suffix         string // Суффикс для запроса (например, ON CONFLICT, RETURNING)
	IgnoreConflict bool   // Флаг игнорирования конфликтов
}
