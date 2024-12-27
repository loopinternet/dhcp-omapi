package omapi

import (
	"net"
	"time"
)

type FailoverState int32

const (
	FailoverStateStartup                   FailoverState = 1
	FailoverStateNormal                                  = 2
	FailoverStateCommunicationsInterrupted               = 3
	FailoverStatePartnerDown                             = 4
	FailoverStatePotentialConflict                       = 5
	FailoverStateRecover                                 = 6
	FailoverStatePaused                                  = 7
	FailoverStateShutdown                                = 8
	FailoverStateRecoverDone                             = 9
	FailoverStateResolutionInterrupted                   = 10
	FailoverStateConflictDone                            = 11
	FailoverStateRecoverWait                             = 254
)

func (state FailoverState) toBytes() []byte {
	return int32ToBytes(int32(state))
}

func (state FailoverState) String() (ret string) {
	switch state {
	case FailoverStateStartup:
		ret = "startup"
	case FailoverStateNormal:
		ret = "normal"
	case FailoverStateCommunicationsInterrupted:
		ret = "communications interrupted"
	case FailoverStatePartnerDown:
		ret = "partner down"
	case FailoverStatePotentialConflict:
		ret = "potential conflict"
	case FailoverStateRecover:
		ret = "recover"
	case FailoverStatePaused:
		ret = "paused"
	case FailoverStateShutdown:
		ret = "shutdown"
	case FailoverStateRecoverDone:
		ret = "recover done"
	case FailoverStateResolutionInterrupted:
		ret = "resolution interrupted"
	case FailoverStateConflictDone:
		ret = "conflict done"
	case FailoverStateRecoverWait:
		ret = "recover wait"
	}

	return
}

type FailoverHierarchy int32

const (
	HierarchyPrimary FailoverHierarchy = iota
	HierarchySecondary
)

func (h FailoverHierarchy) String() (ret string) {
	switch h {
	case HierarchyPrimary:
		ret = "primary"
	case HierarchySecondary:
		ret = "secondary"
	}

	return
}

type Failover struct {
	Name                  string
	PartnerAddress        net.IP
	LocalAddress          net.IP
	PartnerPort           int32
	LocalPort             int32
	MaxOutstandingUpdates int32
	Mclt                  int32 // TODO maybe find a better name
	LoadBalanceMaxSecs    int32
	LoadBalanceHBA        []byte // TODO what type would this be?
	LocalState            FailoverState
	PartnerState          FailoverState
	LocalStos             time.Time // TODO maybe find a better name
	PartnerStos           time.Time // TODO maybe find a better name
	Hierarchy             FailoverHierarchy
	LastPacketSent        time.Time
	LastTimestampReceived time.Time
	Skew                  int32
	MaxResponseDelay      int32
	CurUnackedUpdates     int32
}
