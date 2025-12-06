package types

type StoreResponse struct {
	Aliases              []string            `json:"aliases"`
	RetailChainName      string              `json:"retailChainName"`
	Software             Software            `json:"software"`
	ActiveItems          string              `json:"activeItems"`
	Type                 string              `json:"type"`
	RetailChainID        string              `json:"retailChainId"`
	Setting              Setting             `json:"setting"`
	Contact              Contact             `json:"contact"`
	StoreName            string              `json:"storeName"`
	ID                   string              `json:"id"`
	StoreDomain          string              `json:"storeDomain"`
	RefreshDataFrequency int                 `json:"refreshDataFrequency"`
	StoreID              string              `json:"storeId"`
	CreationDate         string              `json:"creationDate"`
	Modules              Modules             `json:"modules"`
	Typology             string              `json:"typology"`
	ModificationDate     string              `json:"modificationDate"`
	InstallType          string              `json:"installType"`
	Location             Location            `json:"location"`
	Geolocated           bool                `json:"geolocated"`
	Account              Account             `json:"account"`
	Status               StoreStatus         `json:"status"`
	Informations         string              `json:"informations"`
	Applications         []string            `json:"applications"`
	TransmissionMode     string              `json:"transmissionMode"`
	IsLowFrequencyActive bool                `json:"isLowFrequencyActive"`
	TransmissionSystems  TransmissionSystems `json:"transmissionSystems"`
	LastRefreshData      string              `json:"lastRefreshData"`
	Deployment           Deployment          `json:"deployment"`
	Score                int                 `json:"_score"`
}

type Software struct {
	Setting SoftwareSetting `json:"setting"`
}

type SoftwareSetting struct {
	AutoDeploy                 string      `json:"autoDeploy"`
	Owner                      string      `json:"owner"`
	Country                    string      `json:"country"`
	FileName                   string      `json:"fileName"`
	PreferredDateFormatPattern string      `json:"preferredDateFormatPattern"`
	Store                      string      `json:"store"`
	Version                    string      `json:"version"`
	URL                        string      `json:"url"`
	LabelTypes                 []LabelType `json:"labelTypes"`
	LastUpdate                 string      `json:"lastUpdate"`
	Name                       string      `json:"name"`
	AdditionnalInformation     string      `json:"additionnalInformation"`
	ID                         string      `json:"id"`
	Brand                      string      `json:"brand"`
	ExportInformation          ExportInfo  `json:"exportInformation"`
}

type LabelType struct {
	NaturalOrder           string     `json:"naturalOrder"`
	PageCount              string     `json:"pageCount"`
	ImageName              string     `json:"imageName"`
	Pattern                string     `json:"pattern"`
	TypeName               string     `json:"typeName"`
	Scenarios              []Scenario `json:"scenarios"`
	DefaultOrientation     string     `json:"defaultOrientation"`
	DisplayType            string     `json:"displayType"`
	ScreenTechnology       string     `json:"screenTechnology"`
	TransmissionTechnology string     `json:"transmissionTechnology"`
	ScreenColor            string     `json:"screenColor"`
	Width                  string     `json:"width"`
	Animated               string     `json:"animated"`
	TypeID                 string     `json:"typeId"`
	DPI                    string     `json:"dpi"`
	Height                 string     `json:"height"`
	UnitaryUpdateDuration  string     `json:"unitaryUpdateDuration"`
}

type Scenario struct {
	ScenarioName string `json:"scenarioName"`
	ScenarioID   string `json:"scenarioId"`
	ItemCount    string `json:"itemCount"`
}

type ExportInfo struct {
	Date            string `json:"date"`
	SoftwareName    string `json:"softwareName"`
	User            string `json:"user"`
	SoftwareVersion string `json:"softwareVersion"`
}

type Setting struct {
	Name         string `json:"name"`
	ContentMD5   string `json:"contentMD5"`
	ID           string `json:"id"`
	LastModified string `json:"lastModified"`
}

type Contact struct {
	Function string `json:"_function"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Fax      string `json:"fax"`
	Email    string `json:"email"`
}

type Account struct {
	Number string `json:"number"`
}

type Location struct {
	ZipCode      string `json:"zipCode"`
	Country      string `json:"country"`
	StreetNumber string `json:"streetNumber"`
	City         string `json:"city"`
	Street       string `json:"street"`
	State        string `json:"state"`
	GPS          GPS    `json:"gps"`
}

type GPS struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Modules struct {
	Vlink           ModulePlan        `json:"vlink"`
	Adshelf         ModulePlan        `json:"adshelf"`
	AssetManagement ModulePlan        `json:"assetManagement"`
	Geolocation     GeolocationModule `json:"geolocation"`
}

type ModulePlan struct {
	Plan   string `json:"plan"`
	Status string `json:"status"`
}

type GeolocationModule struct {
	Plan     string `json:"plan"`
	Status   string `json:"status"`
	JeegyURL string `json:"jeegyUrl"`
}

type StoreStatus struct {
	Operational  string         `json:"operational"`
	Installation string         `json:"installation"`
	Messages     StatusMessages `json:"messages"`
}

type StatusMessages struct {
	MatchingDaily    any `json:"matchingDaily"`
	Drivres          any `json:"drivres"`
	IntegrationDaily any `json:"integrationDaily"`
	Matching         any `json:"matching"`
	Transmitters     any `json:"transmitters"`
	Integration      any `json:"integration"`
	Labels           any `json:"labels"`
}

type TransmissionSystems struct {
	HighFrequency HighFrequencyTransmission `json:"highFrequency"`
}

type HighFrequencyTransmission struct {
	ModificationDate                 string         `json:"modificationDate"`
	TransmittersEndPoint             string         `json:"transmittersEndPoint"`
	Transmitters                     []Transmitter  `json:"transmitters"`
	TransmittersAuthEndPoint         string         `json:"transmittersAuthEndPoint"`
	Partition                        string         `json:"partition"`
	StoreID                          string         `json:"storeId"`
	Key                              string         `json:"key"`
	LastTransmittersModificationDate string         `json:"lastTransmittersModificationDate"`
	Type                             string         `json:"type"`
	OfflineLabels                    int            `json:"offlineLabels"`
	Connectivity                     HFConnectivity `json:"connectivity"`
	OnlineLabels                     int            `json:"onlineLabels"`
	Status                           string         `json:"status"`
	Flash                            Flash          `json:"flash"`
}

type HFConnectivity struct {
	LastOnlineDate string `json:"lastOnlineDate"`
	Status         string `json:"status"`
}

type Flash struct {
	LastFlashDate string `json:"lastFlashDate"`
}

type Transmitter struct {
	SerialNumber           string                  `json:"serialNumber"`
	TransmitterVersion     string                  `json:"transmitterVersion"`
	Channel                string                  `json:"channel"`
	CreationDate           string                  `json:"creationDate"`
	Provisionning          Provisionning           `json:"provisionning"`
	Version                string                  `json:"version"`
	URL                    string                  `json:"url"`
	ModificationDate       string                  `json:"modificationDate"`
	MacAddress             string                  `json:"macAddress"`
	TransmissionTechnology string                  `json:"transmissionTechnology"`
	ManagedChannel         string                  `json:"managedChannel"`
	Connectivity           TransmitterConnectivity `json:"connectivity"`
	ChannelMode            string                  `json:"channelMode"`
	Name                   string                  `json:"name"`
	InternalVersion        string                  `json:"internalVersion"`
	ID                     string                  `json:"id"`
	TransmitterFamilly     string                  `json:"transmitterFamilly"`
	Status                 string                  `json:"status"`
}

type Provisionning struct {
	DNS                string `json:"dns"`
	Channel            int    `json:"channel"`
	Eth0Address        string `json:"eth0Address"`
	CloudPoint         bool   `json:"cloudPoint"`
	Ble                bool   `json:"ble"`
	JwtSecurityEnabled bool   `json:"jwtSecurityEnabled"`
	Gateway            string `json:"gateway"`
	Eth0Dhcp           bool   `json:"eth0Dhcp"`
	Eth0Netmask        string `json:"eth0Netmask"`
}

type TransmitterConnectivity struct {
	LastOfflineDate string `json:"lastOfflineDate"`
	Status          string `json:"status"`
	LastOnlineDate  string `json:"lastOnlineDate"`
}

type Deployment struct {
	Effective            DeploymentInfo `json:"effective"`
	FirstIntegrationDate string         `json:"firstIntegrationDate"`
	ModificationDate     string         `json:"modificationDate"`
	Expected             DeploymentInfo `json:"expected"`
	FirstSettingDate     string         `json:"firstSettingDate"`
	InstallationRate     int            `json:"installationRate"`
	PrinterInstallation  string         `json:"printerInstallation"`
	Status               string         `json:"status"`
}

type DeploymentInfo struct {
	Accesspoints int     `json:"accesspoints"`
	Printers     int     `json:"printers"`
	EndDate      *string `json:"endDate"`
	Rails        int     `json:"rails"`
	StartDate    *string `json:"startDate"`
	Labels       int     `json:"labels"`
}
