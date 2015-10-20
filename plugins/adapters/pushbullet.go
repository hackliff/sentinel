package adapters

import (
	"fmt"

	"github.com/azer/logger"
	"github.com/xconstruct/go-pushbullet"
)

type PushbulletRadio struct {
	client *pushbullet.Client
	logger *logger.Logger
}

func NewPushbulletRadio(apiKey string) *PushbulletRadio {
	return &PushbulletRadio{
		client: pushbullet.New(apiKey),
		logger: logger.New("sentinel.adapters.pushbullet"),
	}
}

func (p *PushbulletRadio) lookupDevice(name string) (*pushbullet.Device, error) {
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

func (p PushbulletRadio) Send(envelope_ Envelope, message string) error {
	device, err := p.lookupDevice(envelope_.Recipient)
	if err != nil {
		return err
	}
	p.logger.Info("found device iden %s: %s", envelope_.Recipient, device.Iden)
	p.logger.Info("pushing note: %s", envelope_.Title)
	return p.client.PushNote(device.Iden, envelope_.Title, message)
}
