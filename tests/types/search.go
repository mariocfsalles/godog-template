package types

type StoreSearchRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Search   string `json:"search"`
}

type StoreSearchResponse struct {
	Query  Query              `json:"query"`
	Count  int                `json:"count"`
	Values []StoreSearchValue `json:"values"`
}

type StoreSearchValue struct {
	StoreID         string        `json:"storeId"`
	RetailChainName string        `json:"retailChainName"`
	RetailChainID   string        `json:"retailChainId"`
	StoreName       string        `json:"storeName"`
	Location        StoreLocation `json:"location"`
}

type StoreLocation struct {
	ZipCode      string   `json:"zipCode"`
	Country      string   `json:"country"`
	StreetNumber string   `json:"streetNumber"`
	City         string   `json:"city"`
	Street       string   `json:"street"`
	State        string   `json:"state"`
	GPS          StoreGPS `json:"gps"`
}

type StoreGPS struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
