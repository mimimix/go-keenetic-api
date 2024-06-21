package poller

import (
	"keenetic"
	"time"
)

type PollEvent struct {
	IsOnline bool
	Device   *keenetic.Device
}

type Poller struct {
	router       *keenetic.Keenetic
	Interval     int
	Channel      chan *PollEvent
	Devices      map[string]*keenetic.Device
	ticker       *time.Ticker
	isPolling    bool
	isFirstCheck bool
}

func NewPoller(zyxelRouter *keenetic.Keenetic, interval int) *Poller {
	return &Poller{
		router:   zyxelRouter,
		Interval: interval,
		Channel:  make(chan *PollEvent),
		Devices:  make(map[string]*keenetic.Device),
	}
}

func (poller *Poller) CheckDevices() {
	devices, err := poller.router.DeviceList()
	if err != nil {
		return
	}

	for _, device := range *devices {
		_, isExists := poller.Devices[device.Mac]
		if !isExists && !poller.isFirstCheck {
			poller.Channel <- &PollEvent{true, &device}
		} else {
			if poller.Devices[device.Mac].Active != device.Active {
				poller.Channel <- &PollEvent{device.Active, &device}
			}
		}
		poller.Devices[device.Mac] = &device
	}
	poller.isFirstCheck = false
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
