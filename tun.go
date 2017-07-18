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

// Payload get the payload of TUNPacket
func (p TUNPacket) Payload() ([]byte, error) {
	if len(p) < 4 {
		return nil, ErrTooShort
	}
	return p[4:], nil
}
