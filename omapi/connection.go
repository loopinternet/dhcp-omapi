package omapi

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net"
)

type Connection struct {
	authenticator Authenticator
	connection    net.Conn
	inBuffer      *bytes.Buffer
}

// Dial establishes a connection to an OMAPI-enabled server.
func Dial(addr, username, key string) (*Connection, error) {
	con := &Connection{
		authenticator: new(nullAuthenticator),
		inBuffer:      new(bytes.Buffer),
	}

	var newAuth Authenticator = new(nullAuthenticator)

	if len(username) > 0 && len(key) > 0 {
		decodedKey, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			panic(err)
		}
		newAuth = &hmacMD5Authenticator{username, decodedKey, -1}
	}

	tcpConn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	con.connection = tcpConn

	con.sendProtocolInitialization()
	if err := con.receiveProtocolInitialization(); err != nil {
		return nil, err
	}

	if err := con.initializeAuthenticator(newAuth); err != nil {
		return nil, err
	}

	return con, nil
}

func (con *Connection) initializeAuthenticator(auth Authenticator) error {
	if _, ok := auth.(*nullAuthenticator); ok {
		return nil
	}

	message := NewOpenMessage("authenticator")
	for key, value := range auth.AuthObject() {
		message.Object[key] = value
	}

	response, _ := con.Query(message)

	if response.Opcode != OpUpdate {
		return errors.New("received non-update response for open")
	}

	if response.Handle == 0 {
		return errors.New("received invalid authid from server")
	}

	auth.SetAuthID(response.Handle)
	con.authenticator = auth

	return nil
}

// Query sends a message to the server and waits for a reply. It
// returns the underlying response as well as its representation as a
// status. If the message didn't contain a status, the success status
// will be returned instead.
func (con *Connection) Query(msg *Message) (*Message, Status) {
	msg.Sign(con.authenticator)
	con.send(msg.Bytes(false))
	response := con.parseMessage()
	if !response.IsResponseTo(msg) {
		panic("received message is not the desired response")
	}

	// TODO check authid

	return response, response.ToStatus()
}

func (con *Connection) send(data []byte) (n int, err error) {
	return con.connection.Write(data)
}

func (con *Connection) sendProtocolInitialization() {
	buf := newBuffer()
	buf.add(int32(100)) // Protocol version
	buf.add(int32(24))  // Header size
	con.send(buf.bytes())
}

func (con *Connection) read() {
	buf := make([]byte, 2048)
	n, err := con.connection.Read(buf)
	if err != nil {
		panic(err)
	}

	con.inBuffer.Write(buf[0:n])
}

func (con *Connection) waitForN(n int) {
	for con.inBuffer.Len() < n {
		con.read()
	}
}

func (con *Connection) parseStartupMessage() (version, headerSize int32) {
	con.waitForN(8) // version, headerSize

	binary.Read(con.inBuffer, binary.BigEndian, &version)
	binary.Read(con.inBuffer, binary.BigEndian, &headerSize)

	return
}

func (con *Connection) parseMap() map[string][]byte {
	dict := make(map[string][]byte)

	var (
		keyLength   int16
		valueLength int32
		key         []byte
		value       []byte
	)

	for {
		con.waitForN(2) // key length
		binary.Read(con.inBuffer, binary.BigEndian, &keyLength)
		if keyLength == 0 {
			// end of map
			break
		}

		con.waitForN(int(keyLength)) // key
		key = make([]byte, keyLength)
		con.inBuffer.Read(key)

		con.waitForN(4) // value length
		binary.Read(con.inBuffer, binary.BigEndian, &valueLength)
		con.waitForN(int(valueLength)) // value
		value = make([]byte, valueLength)
		con.inBuffer.Read(value)

		dict[string(key)] = value
	}

	return dict
}

func (con *Connection) parseMessage() *Message {
	message := new(Message)
	con.waitForN(24) // authid + authlen + opcode + handle + tid + rid

	var authlen int32

	binary.Read(con.inBuffer, binary.BigEndian, &message.AuthID)
	binary.Read(con.inBuffer, binary.BigEndian, &authlen)
	binary.Read(con.inBuffer, binary.BigEndian, &message.Opcode)
	binary.Read(con.inBuffer, binary.BigEndian, &message.Handle)
	binary.Read(con.inBuffer, binary.BigEndian, &message.TransactionID)
	binary.Read(con.inBuffer, binary.BigEndian, &message.ResponseID)

	message.Message = con.parseMap()
	message.Object = con.parseMap()

	con.waitForN(int(authlen)) // signature
	message.Signature = make([]byte, authlen)
	con.inBuffer.Read(message.Signature)

	return message
}

func (con *Connection) receiveProtocolInitialization() error {
	version, headerSize := con.parseStartupMessage()
	if version != 100 {
		return errors.New("version mismatch")
	}

	if headerSize != 24 {
		return errors.New("header size mismatch")
	}

	return nil
}

func (con *Connection) FindHost(host Host) (Host, error) {
	message := NewOpenMessage("host")

	message.Object = host.toObject()

	response, status := con.Query(message)
	if response.Opcode == OpUpdate {
		return response.ToHost(), nil
	}

	return Host{}, status
}

func (con *Connection) FindLease(lease Lease) (Lease, error) {
	// - IP works
	// - DHCPClientIdentifier works
	// - State does not, even though documentation claims it does
	// - ClientHostname does not, even though documentation claims it does
	message := NewOpenMessage("lease")

	message.Object = lease.toObject()

	response, status := con.Query(message)
	if response.Opcode == OpUpdate {
		return response.ToLease(), nil
	}

	return Lease{}, status
}

// FindFailover finds a failover-state given its name.
func (con *Connection) FindFailover(name string) (Failover, error) {
	message := NewOpenMessage("failover-state")

	message.Object["name"] = []byte(name)

	response, status := con.Query(message)
	if response.Opcode == OpUpdate {
		return response.ToFailover(), nil
	}

	return Failover{}, status
}

// Delete deletes an object from the server, given its handle.
func (con *Connection) Delete(handle int32) error {
	message := NewMessage()
	message.Opcode = OpDelete
	message.Handle = handle

	_, status := con.Query(message)

	if status.IsError() {
		return status
	}

	return nil
}

// CreateHost creates a new host object on the server. The passed
// argument with its fields populated will be sent to the server and
// saved, if possible. The server's representation of the new host
// will be returned.
//
// The returned object will be incomplete compared to the original
// argument, because OMAPI doesn't transfer all information back to
// us.
//
// Example:
//
//	mac, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
//	host := Host{
//		Name:            "new_host",
//		HardwareAddress: mac,
//		HardwareType:    Ethernet,
//		IP:              net.ParseIP("10.0.0.2"),
//		Statements:      `ddns-hostname "the.hostname";`,
//	}
//
//	newHost, err := connection.CreateHost(host)
//	if err != nil {
//		// Couldn't create the new host, error string will tell us why
//	} else {
//		// We successfuly created a new host. newHost will contain the
//		// OMAPI representation of it, including a handle
//	}
func (con *Connection) CreateHost(host Host) (Host, error) {
	message := NewCreateMessage("host")
	message.Object = host.toObject()

	// The server doesn't currently care about Known

	// if host.Known {
	//	message.Object["known"] = True
	// } else {
	//	message.Object["known"] = False
	// }

	response, status := con.Query(message)

	if status.IsError() {
		return Host{}, status
	}

	return response.ToHost(), nil
}

/* func (con *Connection) Update(handle int32, object Object) error {
	// Notes: Cannot change a host's name

	message := NewUpdateMessage(handle)
	message.Object = object.toUpdateObject()

	// Do not transmit empty fields. And as far as we know, unsetting
	// fields doesn't work, anyway.
	for key := range message.Object {
		if len(message.Object[key]) == 0 {
			delete(message.Object, key)
		}
	}

	// The dhcp client identifier can only be changed if it was empty
	// before, so right now, we don't allow setting it at all.
	delete(message.Object, "dhcp-client-identifier")

	// TODO unset fields that cause nasty crashes
	_, status := con.Query(message)

	if status.IsError() {
		return status
	}

	return nil
} */

func (con *Connection) Shutdown() {
	// open Control object, set state to 2, update object, rejoice
}
