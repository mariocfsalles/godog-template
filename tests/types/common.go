package types

type Query struct {
	Search   string      `json:"search,omitempty"`
	Includes []string    `json:"includes,omitempty"`
	Excludes []string    `json:"excludes,omitempty"`
	Sort     []SortField `json:"sort,omitempty"`
	PageSize int         `json:"pageSize,omitempty"`
	Page     int         `json:"page,omitempty"`
}

type SortField map[string]SortOrder

type SortOrder struct {
	Order string `json:"order"`
}
