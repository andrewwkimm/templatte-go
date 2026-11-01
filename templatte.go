package templatte

import "errors"

// ErrEmptyName is returned when a greeting is requested for an empty name.
var ErrEmptyName = errors.New("templatte: empty name")

// Greeting is a message addressed to one recipient.
type Greeting struct {
	// Recipient is the name the greeting addresses.
	Recipient string
	// Message is the full text of the greeting.
	Message string
}

// Greet returns a greeting for name.
func Greet(name string) (Greeting, error) {
	if name == "" {
		return Greeting{}, ErrEmptyName
	}
	return Greeting{Recipient: name, Message: "Hello, " + name + "!"}, nil
}
