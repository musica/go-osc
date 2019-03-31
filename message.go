package osc

type Message struct {
	address   string
	arguments *Arguments
}

func (msg *Message) internal() {}

func NewMessage() *Message {
	return &Message{}
}

func (msg *Message) Clear() *Message {
	msg.address = ""
	msg.arguments = nil

	return msg
}

func (msg *Message) Address() string {
	return msg.address
}

func (msg *Message) SetAddress(address string) *Message {
	msg.address = address
	return msg
}

func (msg *Message) Arguments() *Arguments {
	return msg.arguments
}

func (msg *Message) SetArguments(args *Arguments) *Message {
	msg.arguments = args
	return msg
}

func (msg *Message) MarshalBinary() ([]byte, error) {
	address := createOSCString(msg.address)

	if msg.arguments == nil {
		return address, nil
	}

	typeTags := createOSCString("," + string(msg.arguments.typeTags))
	dataBinary := msg.arguments.dataBinary

	n := len(address) + len(typeTags) + len(dataBinary)
	messageBinary := make([]byte, n)
	n = 0
	n += copy(messageBinary[n:], address)
	n += copy(messageBinary[n:], typeTags)
	n += copy(messageBinary[n:], dataBinary)
	return messageBinary, nil
}
