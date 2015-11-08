package adapters

import (
	"fmt"

	"github.com/azer/logger"
	"github.com/bndr/gotabulate"
)

var TABLE_HEADER = []string{"title", "recipient", "message"}

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

func (s Shell) Send(envelope_ Envelope, message string) error {
	s.logger.Info("\n")
	// NOTE recipient doesn't make a lot of sense here
	table := gotabulate.Create([][]string{
		[]string{envelope_.Title, envelope_.Recipient, message},
	})
	table.SetHeaders(TABLE_HEADER)
	fmt.Println(table.Render("simple"))
	return nil
}
