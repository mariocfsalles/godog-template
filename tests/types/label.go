package types

type LabelRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Search   string `json:"search"`
}

type LabelResponse struct {
	Query  Query        `json:"query"`
	Count  int          `json:"count"`
	Values []LabelValue `json:"values"`
}

type LabelValue struct {
	ModificationDate      string            `json:"modificationDate"`
	LabelID               string            `json:"labelId"`
	Connectivity          LabelConnectivity `json:"connectivity"`
	Page                  int               `json:"page"`
	StoreID               string            `json:"storeId"`
	CreationDate          string            `json:"creationDate"`
	CurrentPage           string            `json:"currentPage"`
	Status                string            `json:"status"`
	Hardware              LabelHardware     `json:"hardware"`
	RetailChainID         string            `json:"retailChainId"`
	Package               LabelPackage      `json:"package"`
	TemperatureDate       string            `json:"temperatureDate"`
	ModificationType      string            `json:"modificationType"`
	LastJoinTimestamp     string            `json:"lastJoinTimestamp"`
	TransmitterID         string            `json:"transmitterId"`
	Transmission          LabelTransmission `json:"transmission"`
	Encryption            LabelEncryption   `json:"encryption"`
	Temperature           float64           `json:"temperature"`
	Billing               LabelBilling      `json:"billing"`
	DeletionDate          *string           `json:"deletionDate"`
	CorrelationID         *string           `json:"correlationId"`
	Matching              LabelMatching     `json:"matching"`
	Page0                 LabelPageInfo     `json:"page_0"`
	PreviousTransmitterID string            `json:"previousTransmitterId"`
	Location              any               `json:"location"` // null no exemplo; tipa depois se precisar
	InStoreDeviceType     string            `json:"inStoreDeviceType"`
}

type LabelConnectivity struct {
	Status           string `json:"status"`
	RSSI             int    `json:"rssi"`
	ModificationDate string `json:"modificationDate"`
	LQI              int    `json:"lqi"`
	SignalQuality    string `json:"signalQuality"`
	LastOnlineDate   string `json:"lastOnlineDate"`
	PreviousStatus   string `json:"previousStatus"`
	LastOfflineDate  string `json:"lastOfflineDate"`
}

type LabelHardware struct {
	NaturalOrder           string  `json:"naturalOrder"`
	PageCount              int     `json:"pageCount"`
	ImageName              string  `json:"imageName"`
	TypeName               string  `json:"typeName"`
	Pattern                string  `json:"pattern"`
	DefaultOrientation     string  `json:"defaultOrientation"`
	StandAlone             string  `json:"standAlone"`
	DisplayType            string  `json:"displayType"`
	ScreenTechnology       string  `json:"screenTechnology"`
	TransmissionTechnology string  `json:"transmissionTechnology"`
	ScreenColor            string  `json:"screenColor"`
	Width                  string  `json:"width"`
	Animated               string  `json:"animated"`
	TypeID                 string  `json:"typeId"`
	DPI                    string  `json:"dpi"`
	Height                 string  `json:"height"`
	UnitaryUpdateDuration  string  `json:"unitaryUpdateDuration"`
	ExtendedFirmware       string  `json:"extendedFirmware"`
	MicrocontrollerType    string  `json:"microcontrollerType"`
	ProductVariant         string  `json:"productVariant"`
	Firmware               string  `json:"firmware"`
	Battery                string  `json:"battery"`
	DeviceType             string  `json:"deviceType"`
	LabelSubSerial         *string `json:"labelSubSerial"`
	LabelSerial            *string `json:"labelSerial"`
}

type LabelPackage struct {
	RegistrationDate string `json:"registrationDate"`
	PackageID        string `json:"packageId"`
	ClaimProvided    bool   `json:"claimProvided"`
}

type LabelTransmission struct {
	RegistrationDate               string `json:"registrationDate"`
	LastSuccessfulTransmissionDate string `json:"lastSuccessfulTransmissionDate"`
	TransmissionDate               string `json:"transmissionDate"`
}

type LabelEncryption struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type LabelBilling struct {
	InactivateDate string `json:"inactivateDate"`
	Status         string `json:"status"`
	ActivateDate   string `json:"activateDate"`
}

type LabelMatching struct {
	MatchingDate string                `json:"matchingDate"`
	Scenario     LabelMatchingScenario `json:"scenario"`
	Items        []LabelMatchingItem   `json:"items"`
	Extended     LabelMatchingExtended `json:"extended"`
}

type LabelMatchingScenario struct {
	Name                  string `json:"name"`
	AutomaticScenarioID   string `json:"automaticScenarioId"`
	AutomaticScenarioName string `json:"automaticScenarioName"`
	ScenarioID            string `json:"scenarioId"`
}

type LabelMatchingExtended struct {
	Battery string `json:"battery"`
}

type LabelMatchingItem struct {
	References       []string    `json:"references"`
	Custom           LabelCustom `json:"custom"`
	Description      string      `json:"description"`
	MatchedItemID    string      `json:"matchedItemId"`
	StoreID          string      `json:"storeId"`
	Extended         any         `json:"extended"`
	ItemID           string      `json:"itemId"`
	ModificationDate string      `json:"modificationDate"`
	Price            float64     `json:"price"`
	Name             string      `json:"name"`
	ID               string      `json:"id"`
	Brand            string      `json:"brand"`
	Status           string      `json:"status"`
}

// ──────
// page_0
// ──────

type LabelPageInfo struct {
	Current  *LabelPageTransmission `json:"current"`
	Expected any                    `json:"expected"`
}

type LabelPageTransmission struct {
	ErrorKey           int                           `json:"errorKey"`
	Metadata           LabelPageTransmissionMetadata `json:"metadata"`
	TimelineID         string                        `json:"timelineId"`
	ErrorMessage       string                        `json:"errorMessage"`
	ExternalID         string                        `json:"externalId"`
	EventType          string                        `json:"eventType"`
	CreationDate       string                        `json:"creationDate"`
	Message            string                        `json:"message"`
	Duration           int                           `json:"duration"`
	TimelineType       string                        `json:"timelineType"`
	ModificationDate   string                        `json:"modificationDate"`
	TransmissionStatus string                        `json:"transmissionStatus"`
	CorrelationID      string                        `json:"correlationId"`
	RetryStrategy      string                        `json:"retryStrategy"`
	Status             string                        `json:"status"`
}

type LabelPageTransmissionMetadata struct {
	OriginalCorrelationID   string  `json:"originalCorrelationId"`
	Retryable               *bool   `json:"retryable"`
	LastRetryDate           *string `json:"lastRetryDate"`
	OriginalExternalID      string  `json:"originalExternalId"`
	ReplacedBy              *string `json:"replacedBy"`
	TransmissionEndDate     string  `json:"transmissionEndDate"`
	ExtClientID             *string `json:"extClientId"`
	ExtCorrelationID        *string `json:"extCorrelationId"`
	OriginalEnqueueDate     string  `json:"originalEnqueueDate"`
	RetryLeft               int     `json:"retryLeft"`
	VTransmitEnqueueDate    string  `json:"vtransmitEnqueueDate"`
	TransmissionEnqueueDate string  `json:"transmissionEnqueueDate"`
	ExternalIDHF            string  `json:"externalIdHf"`
	OriginalEventType       string  `json:"originalEventType"`
	Retry                   *bool   `json:"retry"`
}

// ──────────
// LabelCusto
// ──────────

type LabelCustom struct {
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
	Class                 string `json:"class"`
	EspecialEAN           string `json:"especialean"`
	PromoDesc             string `json:"PROMO_DESC"`
	PrecoUniWas           string `json:"preciouniwas"`
	PVPVatRate            string `json:"PVP_VAT_RATE"`
	Dept                  string `json:"dept"`
	Fechafinal            string `json:"fechafinal"`
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
	ValidadePromocao      string `json:"validade_promocao"`
	SendDate              string `json:"sendDate"`
	PrimUpc               string `json:"primUpc"`
}
