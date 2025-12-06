package snapshot

var Registry = map[string]Config{
	"label_filtered":           NewConfig("label_filtered.json", NormalizeLabelResponse),
	"label_9C1110EB":           NewConfig("label_9C1110EB.json", NormalizeLabelValue),
	"product_filtered":         NewConfig("product_filtered.json", NormalizeProductResponse),
	"storeId_bomdia_pt.009648": NewConfig("storeId_bomdia_pt.009648.json", NormalizeStoreSearchResponse),
	"product_sku_2391674":      NewConfig("product_sku_2391674.json", NormalizeProductValue),
	"search_000206":            NewConfig("search_000206.json", NormalizeStoreResponse),
}
