package pkt

import (
	"errors"
	"io"
	"net"
)

// IPType enum for IP type
type IPType int

const (
	_ IPType = iota
	// SourceIP ip type for source
	SourceIP
	// DestinationIP ip type for source
	DestinationIP
)

// IPPacket represents a single IP packet
type IPPacket []byte

// IPv4PacketHeadLen IPv4 Packet minimum length
const IPv4PacketHeadLen = 20

// IPv6PacketHeadLen IPv6 Packet minimum length
const IPv6PacketHeadLen = 40

var (
	// ErrIPPacketBadVersion IPPacket version not supported
	ErrIPPacketBadVersion = errors.New("IPPacket bad version")
)

// Version returns IPPacket version, 4 or 6, 0 for empty IPPacket
func (p IPPacket) Version() int {
	if len(p) == 0 {
		return 0
	}
	return int(p[0] >> 4)
}

// IP get the source IP, nil for invalid IPPacket
func (p IPPacket) IP(t IPType) (net.IP, error) {
	switch p.Version() {
	case 4:
		{
			if len(p) < IPv4PacketHeadLen {
				return nil, ErrTooShort
			}
			ip := make(net.IP, 4)
			if t == SourceIP {
				copy(ip, p[12:16])
			} else {
				copy(ip, p[16:20])
			}
			return ip, nil
		}
	case 6:
		{
			if len(p) < IPv6PacketHeadLen {
				return nil, ErrTooShort
			}
			ip := make(net.IP, 16)
			if t == SourceIP {
				copy(ip, p[8:24])
			} else {
				copy(ip, p[24:40])
			}
			return ip, nil
		}
	default:
		{
			return nil, ErrIPPacketBadVersion
		}
	}
}

// SetIP set the source IP
func (p IPPacket) SetIP(t IPType, ip net.IP) error {
	switch p.Version() {
	case 4:
		{
			if len(p) < IPv4PacketHeadLen {
				return ErrTooShort
			}
			if len(ip) < net.IPv4len {
				return ErrBadFormat
			}
			if t == SourceIP {
				copy(p[12:16], ip[len(ip)-net.IPv4len:])
			} else {
				copy(p[16:20], ip[len(ip)-net.IPv4len:])
			}
			p.GenerateChecksum()
			return nil
		}
	case 6:
		{
			if len(p) < IPv6PacketHeadLen {
				return ErrTooShort
			}
			if len(ip) < net.IPv6len {
				return ErrBadFormat
			}

			if t == SourceIP {
				copy(p[8:24], ip[len(ip)-net.IPv6len:])
			} else {
				copy(p[24:40], ip[len(ip)-net.IPv6len:])
			}
			return nil
		}
	default:
		{
			return ErrIPPacketBadVersion
		}
	}
}

// GenerateChecksum generate checksum, IPv4 only
func (p IPPacket) GenerateChecksum() error {
	if p.Version() == 4 {
		if len(p) < IPv4PacketHeadLen {
			return ErrTooShort
		}
		sum := uint32(0)
		for i := 0; i < IPv4PacketHeadLen; i += 2 {
			if i != 10 {
				sum += uint32(p[i])<<8 + uint32(p[i+1])
			}
		}
		chksum := ^(uint16(sum) + uint16(sum>>16))
		p[10] = byte(chksum >> 8)
		p[11] = byte(chksum)
	}
	return nil
}

// Length get the length of IPPacket
func (p IPPacket) Length() (int, error) {
	switch p.Version() {
	case 4:
		{
			if len(p) < 4 {
				return -1, ErrTooShort
			}
			return int(p[2])<<4 + int(p[3]), nil
		}
	case 6:
		{
			if len(p) < 6 {
				return -1, ErrTooShort
			}
			return int(p[4])<<4 + int(p[5]) + IPv6PacketHeadLen, nil
		}
	default:
		{
			return -1, ErrIPPacketBadVersion
		}
	}
	return -1, nil
}

// ReadIPPacket read a IPPacket from a io.Reader
func ReadIPPacket(r io.Reader) (IPPacket, error) {
	const HLEN = 6
	// create a minimum header buf, 6 is enough for checking IP version and retrieving IPv4 and IPv6 length
	var err error
	h := make(IPPacket, HLEN, HLEN)
	// read the minimum header
	if _, err := io.ReadFull(r, h); err != nil {
		return nil, err
	}
	// retrieve packet length
	len := 0
	if len, err = h.Length(); err != nil {
		return nil, err
	}
	// append size
	p := make(IPPacket, len)
	copy(p, h)
	// read remaining
	if _, err = io.ReadFull(r, p[HLEN:]); err != nil {
		return nil, err
	}
	return p, nil
}
