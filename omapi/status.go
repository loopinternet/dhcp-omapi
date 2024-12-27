package omapi

type Status struct {
	Code    int32
	Message string
}

var Statuses = []Status{
	Status{0, "success"},
	Status{1, "out of memory"},
	Status{2, "timed out"},
	Status{3, "no available threads"},
	Status{4, "address not available"},
	Status{5, "address in use"},
	Status{6, "permission denied"},
	Status{7, "no pending connections"},
	Status{8, "network unreachable"},
	Status{9, "host unreachable"},
	Status{10, "network down"},
	Status{11, "host down"},
	Status{12, "connection refused"},
	Status{13, "not enough free resources"},
	Status{14, "end of file"},
	Status{15, "socket already bound"},
	Status{16, "task is done"},
	Status{17, "lock busy"},
	Status{18, "already exists"},
	Status{19, "ran out of space"},
	Status{20, "operation canceled"},
	Status{21, "sending events is not allowed"},
	Status{22, "shutting down"},
	Status{23, "not found"},
	Status{24, "unexpected end of input"},
	Status{25, "failure"},
	Status{26, "I/O error"},
	Status{27, "not implemented"},
	Status{28, "unbalanced parentheses"},
	Status{29, "no more"},
	Status{30, "invalid file"},
	Status{31, "bad base64 encoding"},
	Status{32, "unexpected token"},
	Status{33, "quota reached"},
	Status{34, "unexpected error"},
	Status{35, "already running"},
	Status{36, "host unknown"},
	Status{37, "protocol version mismatch"},
	Status{38, "protocol error"},
	Status{39, "invalid argument"},
	Status{40, "not connected"},
	Status{41, "data not yet available"},
	Status{42, "object unchanged"},
	Status{43, "more than one object matches key"},
	Status{44, "key conflict"},
	Status{45, "parse error(s) occurred"},
	Status{46, "no key specified"},
	Status{47, "zone TSIG key not known"},
	Status{48, "invalid TSIG key"},
	Status{49, "operation in progress"},
	Status{50, "DNS format error"},
	Status{51, "DNS server failed"},
	Status{52, "no such domain"},
	Status{53, "not implemented"},
	Status{54, "refused"},
	Status{55, "domain already exists"},
	Status{56, "RRset already exists"},
	Status{57, "no such RRset"},
	Status{58, "not authorized"},
	Status{59, "not a zone"},
	Status{60, "bad DNS signature"},
	Status{61, "bad DNS key"},
	Status{62, "clock skew too great"},
	Status{63, "no root zone"},
	Status{64, "destination address required"},
	Status{65, "cross-zone update"},
	Status{66, "no TSIG signature"},
	Status{67, "not equal"},
	Status{68, "connection reset by peer"},
	Status{69, "unknown attribute"},
}

// IsError returns true if the status is describing an error.
func (s Status) IsError() bool {
	return s.Code > 0
}

func (s Status) Error() string {
	return s.Message
}
