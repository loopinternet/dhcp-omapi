package omapi

import (
	"bytes"
	"net"
	"time"
)

type Message struct {
	AuthID        int32
	Opcode        Opcode
	Handle        int32
	TransactionID int32
	ResponseID    int32
	Message       map[string][]byte
	Object        map[string][]byte
	Signature     []byte
}

func NewMessage() *Message {
	rng.Lock()
	tid := rng.Int31()
	rng.Unlock()

	msg := &Message{
		TransactionID: tid,
		Message:       make(map[string][]byte),
		Object:        make(map[string][]byte),
	}

	return msg
}

func NewOpenMessage(typeName string) *Message {
	message := NewMessage()
	message.Opcode = OpOpen
	message.Message["type"] = []byte(typeName)

	return message
}

func NewCreateMessage(typeName string) *Message {
	message := NewOpenMessage(typeName)
	message.Message["create"] = True
	message.Message["exclusive"] = True

	return message
}

func NewDeleteMessage(handle int32) *Message {
	message := NewMessage()
	message.Opcode = OpDelete
	message.Handle = handle

	return message
}

func (m *Message) Bytes(forSigning bool) []byte {
	ret := newBuffer()
	if !forSigning {
		ret.add(m.AuthID)
	}

	ret.add(int32(len(m.Signature)))
	ret.add(m.Opcode)
	ret.add(m.Handle)
	ret.add(m.TransactionID)
	ret.add(m.ResponseID)
	ret.addMap(m.Message)
	ret.addMap(m.Object)
	if !forSigning {
		ret.add(m.Signature)
	}

	return ret.buffer.Bytes()
}

func (m *Message) Sign(auth Authenticator) {
	m.AuthID = auth.AuthID()
	m.Signature = auth.Sign(m)
}

func (m *Message) Verify(auth Authenticator) bool {
	return bytes.Equal(auth.Sign(m), m.Signature)
}

func (m *Message) IsResponseTo(other *Message) bool {
	return m.ResponseID == other.TransactionID
}

func (m *Message) ToHost() Host {
	return Host{
		Name:                 string(m.Object["name"]),
		HardwareAddress:      net.HardwareAddr(m.Object["hardware-address"]),
		HardwareType:         HardwareType(bytesToInt32(m.Object["hardware-type"])),
		DHCPClientIdentifier: m.Object["dhcp-client-identifier"],
		IP:                   net.IP(m.Object["ip-address"]),
		Handle:               m.Handle,
	}
}

func (m *Message) ToStatus() Status {
	if m.Opcode != OpStatus {
		return Statuses[0]
	}

	return Statuses[bytesToInt32(m.Message["result"])]
}

func (m *Message) ToLease() Lease {
	state := bytesToInt32(m.Object["state"])
	host := bytesToInt32(m.Object["host"])
	ends := bytesToInt32(m.Object["ends"])
	tstp := bytesToInt32(m.Object["tstp"])
	atsfp := bytesToInt32(m.Object["atsfp"])
	cltt := bytesToInt32(m.Object["cltt"])

	return Lease{
		State:                LeaseState(state),
		IP:                   net.IP(m.Object["ip-address"]),
		DHCPClientIdentifier: m.Object["dhcp-client-identifier"],
		ClientHostname:       string(m.Object["client-hostname"]),
		Host:                 host,
		HardwareAddress:      net.HardwareAddr(m.Object["hardware-address"]),
		HardwareType:         HardwareType(bytesToInt32(m.Object["hardware-type"])),
		Ends:                 time.Unix(int64(ends), 0),
		Tstp:                 time.Unix(int64(tstp), 0),
		Atsfp:                time.Unix(int64(atsfp), 0),
		Cltt:                 time.Unix(int64(cltt), 0),
		Handle:               m.Handle,
	}
}

func (m *Message) ToFailover() Failover {
	partnerPort := bytesToInt32(m.Object["partner-port"])
	localPort := bytesToInt32(m.Object["local-port"])
	maxOutstandingUpdates := bytesToInt32(m.Object["max-outstanding-updates"])
	mclt := bytesToInt32(m.Object["mclt"])
	loadBalanceMaxSecs := bytesToInt32(m.Object["load-balance-max-secs"])
	localState := bytesToInt32(m.Object["local-state"])
	partnerState := bytesToInt32(m.Object["partner-state"])
	localStos := bytesToInt32(m.Object["local-stos"])
	partnerStos := bytesToInt32(m.Object["partner-stos"])
	hierarchy := bytesToInt32(m.Object["hierarchy"])
	lastPacketSent := bytesToInt32(m.Object["last-packet-sent"])
	lastTimestampReceived := bytesToInt32(m.Object["last-timestamp-received"])
	skew := bytesToInt32(m.Object["skew"])
	maxResponseDelay := bytesToInt32(m.Object["max-response-delay"])
	curUnackedUpdates := bytesToInt32(m.Object["cur-unacked-updates"])

	return Failover{
		Name:                  string(m.Object["name"]),
		PartnerAddress:        net.IP(m.Object["partner-address"]),
		LocalAddress:          net.IP(m.Object["local-address"]),
		PartnerPort:           partnerPort,
		LocalPort:             localPort,
		MaxOutstandingUpdates: maxOutstandingUpdates,
		Mclt:                  mclt,
		LoadBalanceMaxSecs:    loadBalanceMaxSecs,
		LoadBalanceHBA:        m.Object["load-balance-hba"],
		LocalState:            FailoverState(localState),
		PartnerState:          FailoverState(partnerState),
		LocalStos:             time.Unix(int64(localStos), 0),
		PartnerStos:           time.Unix(int64(partnerStos), 0),
		Hierarchy:             FailoverHierarchy(hierarchy),
		LastPacketSent:        time.Unix(int64(lastPacketSent), 0),
		LastTimestampReceived: time.Unix(int64(lastTimestampReceived), 0),
		Skew:                  skew,
		MaxResponseDelay:      maxResponseDelay,
		CurUnackedUpdates:     curUnackedUpdates,
	}
}
