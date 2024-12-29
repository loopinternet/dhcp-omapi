package omapi

import (
	"encoding/json"
	"net"
	"time"
)

type LeaseState int32

const (
	_              = iota
	LeaseStateFree = iota
	LeaseStateActive
	LeaseStateExpired
	LeaseStateReleased
	LeaseStateAbandoned
	LeaseStateReset
	LeaseStateBackup
	LeaseStateReserved
	LeaseStateBootp
)

func (state LeaseState) String() (ret string) {
	switch state {
	case LeaseStateFree:
		ret = "free"
	case LeaseStateActive:
		ret = "active"
	case LeaseStateExpired:
		ret = "expired"
	case LeaseStateReleased:
		ret = "released"
	case LeaseStateAbandoned:
		ret = "abandoned"
	case LeaseStateReset:
		ret = "reset"
	case LeaseStateBackup:
		ret = "backup"
	case LeaseStateReserved:
		ret = "reserved"
	case LeaseStateBootp:
		ret = "bootp"
	}

	return
}

func (state LeaseState) toBytes() []byte {
	return int32ToBytes(int32(state))
}

func (state LeaseState) MarshalText() ([]byte, error) {
	return []byte(state.String()), nil
}

type HardwareType int32

const (
	Ethernet  HardwareType = 1
	TokenRing              = 6
	FDDI                   = 8
)

func (hw HardwareType) toBytes() []byte {
	return int32ToBytes(int32(hw))
}

func (hw HardwareType) String() (ret string) {
	switch hw {
	case Ethernet:
		ret = "Ethernet"
	case TokenRing:
		ret = "Token ring"
	case FDDI:
		ret = "FDDI"
	}

	return
}

func (hw HardwareType) MarshalText() ([]byte, error) {
	return []byte(hw.String()), nil
}

type Lease struct {
	State                LeaseState
	IP                   net.IP
	DHCPClientIdentifier []byte
	ClientHostname       string
	Host                 int32 // TODO figure out what to do with handles
	// Subnet, Pool, BillingClass are "currently not supported" by the dhcpd
	HardwareAddress net.HardwareAddr
	HardwareType    HardwareType
	Ends            time.Time
	// TODO maybe find nicer names for these times
	Tstp   time.Time
	Atsfp  time.Time
	Cltt   time.Time
	Handle int32
}

func (lease Lease) String() string {

	tmp := struct {
		State                LeaseState   `json:"state"`
		IP                   net.IP       `json:"ip"`
		DHCPClientIdentifier []byte       `json:"dhcp-client-identifier"`
		ClientHostname       string       `json:"client-hostname"`
		Host                 int32        `json:"host"`
		HardwareAddress      string       `json:"hardware-address"`
		HardwareType         HardwareType `json:"hardware-type"`
		Ends                 time.Time    `json:"lease-end-time"`
		Tstp                 time.Time    `json:"tstp"`
		Atsfp                time.Time    `json:"atsfp"`
		Cltt                 time.Time    `json:"cltt"`
		Handle               int32        `json:"handle"`
	}{
		State:                lease.State,
		IP:                   lease.IP,
		DHCPClientIdentifier: lease.DHCPClientIdentifier,
		ClientHostname:       lease.ClientHostname,
		Host:                 lease.Host,
		HardwareAddress:      net.HardwareAddr(lease.HardwareAddress).String(),
		HardwareType:         lease.HardwareType,
		Ends:                 lease.Ends,
		Tstp:                 lease.Tstp,
		Atsfp:                lease.Atsfp,
		Cltt:                 lease.Cltt,
		Handle:               lease.Handle,
	}

	ret, _ := json.Marshal(tmp)

	return string(ret)
}

func (lease Lease) toObject() map[string][]byte {
	object := make(map[string][]byte)

	// TODO check if sending the state in an update will cause an
	// error
	if lease.State > 0 {
		object["state"] = lease.State.toBytes()
	} else {
		object["state"] = nil
	}

	// TODO check if sending the IP in an update will cause an
	// error
	if len([]byte(lease.IP)) > 0 {
		object["ip-address"] = []byte(lease.IP)[12:]
	} else {
		object["ip-address"] = nil
	}

	object["dhcp-client-identifier"] = lease.DHCPClientIdentifier
	object["client-hostname"] = []byte(lease.ClientHostname)
	object["hardware-address"] = []byte(lease.HardwareAddress)

	if lease.HardwareType == 0 {
		object["hardware-type"] = nil
	} else {
		object["hardware-type"] = lease.HardwareType.toBytes()
	}

	return object
}
