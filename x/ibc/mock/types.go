package mock

import (
	"math"

	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	ibc_channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "ibcmock"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

type MockPacket struct {
	Data []byte
}

var _ ibc_channel.Packet = MockPacket{}

func (_ MockPacket) Timeout() uint64 {
	return math.MaxUint64
}

func (_ MockPacket) Commit() []byte {
	return []byte{0}
}

type MockConsensusState struct {
	PreventNil bool `amino:write-empty` // It is neccessary because state/value panic if amino return nil.
}

var _ client.ConsensusState = MockConsensusState{}

func (_ MockConsensusState) Kind() client.Kind {
	return 0
}

func (_ MockConsensusState) GetHeight() uint64 {
	return 1
}

func (_ MockConsensusState) GetRoot() commitment.Root {
	return MockRoot{}
}

func (_ MockConsensusState) Validate(client.Header) (client.ConsensusState, error) {
	return MockConsensusState{}, nil
}

func (_ MockConsensusState) Equivocation(client.Header, client.Header) bool {
	return false
}

type MockRoot struct {
	PreventNil bool `amino:write-empty` // It is neccessary because state/value panic if amino return nil.
}

var _ commitment.Root = MockRoot{}

func (_ MockRoot) CommitmentKind() string {
	return "mock"
}

type MockPath struct {
	PreventNil bool `amino:write-empty` // It is neccessary because state/value panic if amino return nil.
}

var _ commitment.Path = MockPath{}

func (_ MockPath) CommitmentKind() string {
	return "mock"
}

func (_ MockPath) Pathify([]byte) []byte {
	return []byte{}
}

type MockProof struct {
	PreventNil bool `amino:write-empty` // It is neccessary because state/value panic if amino return nil.
}

var _ commitment.Proof = MockProof{}

func (_ MockProof) CommitmentKind() string {
	return "mock"
}

func (_ MockProof) GetKey() []byte {
	return []byte("mock")
}

func (_ MockProof) Verify(commitment.Root, commitment.Path, []byte) error {
	return nil
}
