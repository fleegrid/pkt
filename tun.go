package pkt

// TUNPacket aliase to []byte
type TUNPacket []byte

// Flags get the flags
func (p TUNPacket) Flags() (uint16, error) {
	if len(p) < 2 {
		return 0, ErrTooShort
	}
	return uint16(p[0])<<8 + uint16(p[1]), nil
}

// Proto get the proto of TUN device
func (p TUNPacket) Proto() (uint16, error) {
	if len(p) < 4 {
		return 0, ErrTooShort
	}
	return uint16(p[2])<<8 + uint16(p[3]), nil
}

// SetProto set the proto of TUN device
func (p TUNPacket) SetProto(proto uint16) error {
	if len(p) < 4 {
		return ErrTooShort
	}
	p[2] = byte(proto >> 8)
	p[3] = byte(proto)
	return nil
}

// CopyPayload copy bytes to payload
func (p TUNPacket) CopyPayload(b []byte) (int, error) {
	if len(p) < 4 {
		return 0, ErrTooShort
	}
	return copy(p[4:], b), nil
}

// Payload get the payload of TUNPacket
func (p TUNPacket) Payload() ([]byte, error) {
	if len(p) < 4 {
		return nil, ErrTooShort
	}
	return p[4:], nil
}
