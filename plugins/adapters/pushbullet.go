package adapters

import (
	"fmt"

	"github.com/azer/logger"
	"github.com/xconstruct/go-pushbullet"
)

func init() {
	log.Info("registering adapter: pushbullet")
	Add("pushbullet", func(conf map[string]string) (Plugin, error) {
		// NOTE read it from the environment ?
		apiKey, ok := conf["api-key"]
		if !ok {
			return nil, fmt.Errorf("no api key provided for pushbullet adapter")
		}
		return NewPushbullet(apiKey), nil
	})
}

type Pushbullet struct {
	client *pushbullet.Client
	logger *logger.Logger
}

func NewPushbullet(apiKey string) *Pushbullet {
	return &Pushbullet{
		client: pushbullet.New(apiKey),
		logger: logger.New("sentinel.adapters.pushbullet"),
	}
}

func (p *Pushbullet) lookupDevice(name string) (*pushbullet.Device, error) {
	devs, err := p.client.Devices()
	if err != nil {
		return nil, err
	}

	for _, device := range devs {
		if device.Nickname == name {
			return device, nil
		}
	}
	return nil, fmt.Errorf("device %s not found\n", name)
}

func (p Pushbullet) Send(envelope_ Envelope, message string) error {
	device, err := p.lookupDevice(envelope_.Recipient)
	if err != nil {
		return err
	}
	p.logger.Info("found device iden %s: %s", envelope_.Recipient, device.Iden)
	p.logger.Info("pushing note: %s", envelope_.Title)
	return p.client.PushNote(device.Iden, envelope_.Title, message)
}
