package omapi

type Opcode int32

const (
	_             = iota
	OpOpen Opcode = iota
	OpRefresh
	OpUpdate
	OpNotify
	OpStatus
	OpDelete
)

func (opcode Opcode) String() (ret string) {
	switch opcode {
	case OpOpen:
		ret = "open"
	case OpRefresh:
		ret = "refresh"
	case OpUpdate:
		ret = "update"
	case OpNotify:
		ret = "notify"
	case OpStatus:
		ret = "status"
	case OpDelete:
		ret = "delete"
	}

	return
}
