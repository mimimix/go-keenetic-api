package keenetic

type Device struct {
	Mac       string `json:"mac"`
	Via       string `json:"via"`
	Ip        string `json:"ip"`
	Hostname  string `json:"hostname"`
	Name      string `json:"name"`
	Interface struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"interface"`
	Expires         int64  `json:"expires"`
	Registered      bool   `json:"registered"`
	Access          string `json:"access"`
	Schedule        string `json:"schedule"`
	Active          bool   `json:"active"`
	Rxbytes         int64  `json:"rxbytes"`
	Txbytes         int64  `json:"txbytes"`
	Uptime          int64  `json:"uptime"`
	FirstSeen       int64  `json:"first-seen"`
	LastSeen        int64  `json:"last-seen"`
	Link            string `json:"link"`
	AutoNegotiation bool   `json:"auto-negotiation"`
	Speed           int64  `json:"speed"`
	Duplex          bool   `json:"duplex"`
	EverSeen        bool   `json:"ever-seen"`
	TrafficShape    struct {
		Rx       int64  `json:"rx"`
		Tx       int64  `json:"tx"`
		Mode     string `json:"mode"`
		Schedule string `json:"schedule"`
	} `json:"traffic-shape"`
}

type deviceListResponse struct {
	Show struct {
		Ip struct {
			Hotspot struct {
				Host []Device `json:"host"`
			} `json:"hotspot"`
		} `json:"ip"`
	} `json:"show"`
}

func (zyxel *Keenetic) DeviceList() (*[]Device, error) {
	var result deviceListResponse
	_, err := zyxel.Request.R().
		SetBody(`{"show":{"ip":{"hotspot":{}}}}`).
		SetSuccessResult(&result).
		Post("/rci/")
	if err != nil {
		return nil, err
	}
	return &result.Show.Ip.Hotspot.Host, nil
}
