package types

type OrderRequest struct {
	Customer string   `json:"customer,omitempty"`
	Items    []string `json:"items,omitempty"`
	Status   string   `json:"status,omitempty"`
}

type OrderResponse struct {
	ID        string   `json:"id,omitempty"`
	Customer  string   `json:"customer,omitempty"`
	Items     []string `json:"items,omitempty"`
	Status    string   `json:"status,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	UpdatedAt string   `json:"updatedAt,omitempty"`
}
