package keenetic

import (
	"slices"
	"time"
)

type PollEvent struct {
	IsOnline bool
	Device   *Device
}

type Poller struct {
	router       *Keenetic
	Interval     int64
	Channel      chan *PollEvent
	Devices      map[string]*Device
	ticker       *time.Ticker
	isPolling    bool
	isFirstCheck bool
	LastOnline   LastOnline
}

type LastOnline map[string]time.Time

func NewPoller(zyxelRouter *Keenetic, interval int64) *Poller {
	return &Poller{
		router:     zyxelRouter,
		Interval:   interval,
		Channel:    make(chan *PollEvent),
		Devices:    make(map[string]*Device),
		LastOnline: make(LastOnline),
	}
}

func (poller *Poller) CheckDevices() {
	devices, err := poller.router.DeviceList()
	if err != nil {
		return
	}

	nowMacDevices := make([]string, len(*devices))
	for _, device := range *devices {
		nowMacDevices = append(nowMacDevices, device.Mac)
		_, isExists := poller.Devices[device.Mac]
		if !isExists && !poller.isFirstCheck {
			poller.Channel <- &PollEvent{device.Active, &device}
		} else {
			if poller.Devices[device.Mac].Active != device.Active {
				poller.Channel <- &PollEvent{device.Active, &device}
				if !device.Active {
					poller.LastOnline[device.Mac] = time.Now()
				} else {
					delete(poller.Devices, device.Mac)
				}
			}
		}
		poller.Devices[device.Mac] = &device
	}
	for mac, device := range poller.Devices {
		if !slices.Contains(nowMacDevices, mac) {
			poller.Channel <- &PollEvent{false, device}
			poller.LastOnline[device.Mac] = time.Now()
		} else {
			delete(poller.Devices, mac)
		}
	}
	poller.isFirstCheck = false
}

func (poller *Poller) GetLastOnline(mac string) (time.Time, bool) {
	t, isExists := poller.LastOnline[mac]
	return t, isExists
}

func (poller *Poller) Start() {
	if poller.isPolling {
		return
	}
	ticker := time.NewTicker(time.Duration(poller.Interval) * time.Second)
	poller.ticker = ticker
	poller.isPolling = true
	go func() {
		for {
			<-ticker.C
			poller.CheckDevices()
		}
	}()
}
func (poller *Poller) Stop() {
	poller.ticker.Stop()
	poller.isPolling = false
}
