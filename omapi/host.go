package omapi

import "net"

type Host struct {
	Name                 string
	Group                int32 // TODO
	HardwareAddress      net.HardwareAddr
	HardwareType         HardwareType
	DHCPClientIdentifier []byte
	IP                   net.IP
	Statements           string // Not populated by OMAPI
	Known                bool   // Not populated by OMAPI
	Handle               int32
}

func (host Host) toObject() map[string][]byte {
	object := make(map[string][]byte)

	object["name"] = []byte(host.Name)

	if len([]byte(host.IP)) > 0 {
		object["ip-address"] = []byte(host.IP.To4())
	} else {
		object["ip-address"] = nil
	}

	object["hardware-address"] = []byte(host.HardwareAddress)
	if host.HardwareType == 0 {
		object["hardware-type"] = nil
	} else {
		object["hardware-type"] = host.HardwareType.toBytes()
	}

	// TODO remove statements field when updating an object, to work around bug
	object["statements"] = []byte(host.Statements)

	object["dhcp-client-identifier"] = host.DHCPClientIdentifier

	return object
}
