package quickfix

func WithTagOrderBody(order TagOrder) Opt {
	return func(m *Message) {
		m.Body.initWithOrdering(order)
	}
}

func WithTagOrderHeader(order TagOrder) Opt {
	return func(m *Message) {
		m.Header.initWithOrdering(order)
	}
}

func NewMessageBuilder(opts ...Opt) *Message {
	m := newMessage()
	for _, opt := range opts {
		opt(m)
	}

	return m
}

var _ TagOrder = AscFieldOrder

// AscFieldOrder orders tags by their integer value while treating Party related tags
// separately placing them in 453,447,448,452 order at the end of the message's body
func AscFieldOrder(i, j Tag) bool {
	ordering := func(t Tag) uint32 {
		switch t {
		case 453: //NoPartyIds
			return 1
		case 447: //PartyIdSource
			return 2
		case 448: //PartyId
			return 3
		case 452: //PartyRole
			return 4
		}

		// all other tags will precede party
		return 0
	}

	orderi := ordering(i)
	orderj := ordering(j)

	switch {
	case orderi < orderj:
		return true
	case orderi > orderj:
		return false
	}

	return i < j
}
