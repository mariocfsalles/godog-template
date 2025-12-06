package snapshot

import "orchestrator/types"

func NormalizeStoreResponse(s *types.StoreResponse) {
	// top-level voláteis
	s.CreationDate = ""
	s.ModificationDate = ""
	s.LastRefreshData = ""

	// Deployment
	s.Deployment.FirstIntegrationDate = ""
	s.Deployment.ModificationDate = ""
	s.Deployment.FirstSettingDate = ""

	if s.Deployment.Effective.EndDate != nil {
		*s.Deployment.Effective.EndDate = ""
	}
	if s.Deployment.Effective.StartDate != nil {
		*s.Deployment.Effective.StartDate = ""
	}
	if s.Deployment.Expected.EndDate != nil {
		*s.Deployment.Expected.EndDate = ""
	}
	if s.Deployment.Expected.StartDate != nil {
		*s.Deployment.Expected.StartDate = ""
	}
	s.Deployment.Effective.Labels = 0
	s.Deployment.Expected.Labels = 0

	// TransmissionSystems.HighFrequency
	hf := &s.TransmissionSystems.HighFrequency

	hf.ModificationDate = ""
	hf.LastTransmittersModificationDate = ""
	hf.Connectivity.LastOnlineDate = ""

	hf.OfflineLabels = 0
	hf.OnlineLabels = 0

	for i := range hf.Transmitters {
		t := &hf.Transmitters[i]
		t.CreationDate = ""
		t.ModificationDate = ""
		t.Connectivity.LastOnlineDate = ""
		t.Connectivity.LastOfflineDate = ""
		t.Name = "" // IP
	}
}

func NormalizeStoreSearchResponse(s *types.StoreSearchResponse) {}

func NormalizeProductResponse(s *types.ProductResponse) {
	for i := range s.Values {
		v := &s.Values[i]
		v.ModificationDate = ""
		v.Custom.StockInTransitNight = ""
		v.Custom.SOHNight = ""
	}
}

func NormalizeProductValue(s *types.ProductValue) {
	s.ModificationDate = ""
	s.Custom.SOHNight = ""
	s.Custom.StockInTransitNight = ""
}

func NormalizeLabelResponse(s *types.LabelResponse) {
	for i := range s.Values {
		v := &s.Values[i]
		v.ModificationDate = ""
		v.Connectivity.RSSI = 0
		v.LastJoinTimestamp = ""
	}
}

func NormalizeLabelValue(s *types.LabelValue) {
	s.ModificationDate = ""
	s.TemperatureDate = ""
	s.LastJoinTimestamp = ""
}
