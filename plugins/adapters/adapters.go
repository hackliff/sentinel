// adapters package stores various plugins aimed at bot communication with an
// operator.
package adapters

type Envelope struct {
	Title     string
	Recipient string
}

type Adapter interface {
	Send(Envelope, string) error
}
