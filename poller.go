package keenetic

import (
	"time"
)

type PollEvent struct {
	IsOnline bool
	Device   *Device
}

type Poller struct {
	router       *Keenetic
	Interval     int
	Channel      chan *PollEvent
	Devices      map[string]*Device
	ticker       *time.Ticker
	isPolling    bool
	isFirstCheck bool
}

func NewPoller(zyxelRouter *Keenetic, interval int) *Poller {
	return &Poller{
		router:   zyxelRouter,
		Interval: interval,
		Channel:  make(chan *PollEvent),
		Devices:  make(map[string]*Device),
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
