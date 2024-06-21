package poller

import (
	"time"
	"zyxel/pkg/zyxel"
)

type PollEvent struct {
	IsOnline bool
	Device   *zyxel.Device
}

type Poller struct {
	zyxel        *zyxel.Zyxel
	Interval     int
	Channel      chan *PollEvent
	Devices      map[string]*zyxel.Device
	ticker       *time.Ticker
	isPolling    bool
	isFirstCheck bool
}

func NewPoller(zyxelRouter *zyxel.Zyxel, interval int) *Poller {
	return &Poller{
		zyxel:    zyxelRouter,
		Interval: interval,
		Channel:  make(chan *PollEvent),
		Devices:  make(map[string]*zyxel.Device),
	}
}

func (poller *Poller) CheckDevices() {
	devices, err := poller.zyxel.DeviceList()
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
