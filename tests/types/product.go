package types

type ProductRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Search   string `json:"search"`
}

type ProductResponse struct {
	Query  Query          `json:"query"`
	Count  int            `json:"count"`
	Values []ProductValue `json:"values"`
}

type ProductValue struct {
	References       []string        `json:"references"`
	Custom           ProductCustom   `json:"custom"`
	Description      string          `json:"description"`
	StoreID          string          `json:"storeId"`
	CreationDate     string          `json:"creationDate"`
	Extended         any             `json:"extended"`
	Matching         ProductMatching `json:"matching"`
	RetailChainID    string          `json:"retailChainId"`
	ItemID           string          `json:"itemId"`
	ModificationDate string          `json:"modificationDate"`
	DeletionDate     *string         `json:"deletionDate"`
	Price            float64         `json:"price"`
	Name             string          `json:"name"`
	ID               string          `json:"id"`
	Brand            string          `json:"brand"`
	Status           string          `json:"status"`
}

type ProductMatching struct {
	Count   int      `json:"count"`
	Matched bool     `json:"matched"`
	Labels  []string `json:"labels"`
}

type ProductCustom struct {
	EANPVPWas             string `json:"eanpvpwas"`
	PVRWasDescription     string `json:"PVR_WAS_Description"`
	DiscountDesc          string `json:"DISCOUNT_DESC"`
	CustomSKU             string `json:"customsku"`
	Desconto              string `json:"desconto"`
	DescWas               string `json:"DESC_WAS"`
	Versiesp              string `json:"versiesp"`
	PVRISDescription      string `json:"PVR_IS_Description"`
	TicketSubtype         string `json:"TICKET_SUBTYPE"`
	PVR                   string `json:"PVR"`
	DescIS                string `json:"DESC_IS"`
	Fechainicio           string `json:"fechainicio"`
	SendDate              string `json:"send_date"`
	Class                 string `json:"class"`
	EspecialEAN           string `json:"especialean"`
	PromoDesc             string `json:"PROMO_DESC"`
	PrecoUniWas           string `json:"preciouniwas"`
	PVPVatRate            string `json:"PVP_VAT_RATE"`
	Dept                  string `json:"dept"`
	PrimUPC               string `json:"prim_upc"`
	Fechafinal            string `json:"fechafinal"`
	ValidadePromocao      string `json:"validade_promocao"`
	PVPUOMDesc            string `json:"PVP_UOM_DESC"`
	PrecoUniS             string `json:"preciounis"`
	PromoSimplifiedLayout string `json:"promo_simplified_layout"`
	PrecoAntes            string `json:"precioantes"`
	ShowVatRateStamp      string `json:"SHOW_VAT_RATE_STAMP"`
	Unidades              string `json:"unidades"`
	Subclass              string `json:"subclass"`
	EspecialID            string `json:"especialid"`
	VATRateDesc           string `json:"VAT_RATE_DESC"`
	Status                string `json:"status"`
	QtyOnOrderNight       string `json:"qty_on_order_night"`
	DemoStockNight        string `json:"demo_stock_night"`
	PresStockNight        string `json:"pres_stock_night"`
	StockInTransitNight   string `json:"stock_in_transit_night"`
	RuptAlertNightESL     string `json:"rupt_alert_night_esl"`
	SRPInitNight          string `json:"srp_init_night"`
	TagDateNight          string `json:"tag_date_night"`
	RuptAlertNight        string `json:"rupt_alert_night"`
	QtySold               string `json:"qty_sold"`
	SOHNight              string `json:"soh_night"`
	SendDateNight         string `json:"send_date_night"`
	StatusNight           string `json:"status_night"`
}
