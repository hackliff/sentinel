package adapters

import "github.com/azer/logger"

func init() {
	log.Info("registering adapter: shell")
	Add("shell", func(conf map[string]string) (Plugin, error) {
		return &Shell{
			logger: logger.New("sentinel.plugins.adapters.shell"),
		}, nil
	})
}

type Shell struct {
	logger *logger.Logger
}

func (p Shell) Send(envelope_ Envelope, message string) error {
	p.logger.Info("\ttitle:     %s", envelope_.Title)
	// NOTE recipient doesn't make a lot of sense here
	p.logger.Info("\trecipient: %s", envelope_.Recipient)
	p.logger.Info("\tmessage:   %s", message)
	return nil
}
