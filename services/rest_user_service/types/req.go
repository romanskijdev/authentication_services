package typesm

type Response struct {
	TotalCount *uint64     `json:"total_count,omitempty"`
	Count      int         `json:"count,omitempty"`
	Error      *string     `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}
