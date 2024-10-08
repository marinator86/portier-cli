package messages

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type ConnectionID string

type MessageType string

const (
	// ConnectionOpenMessage is a message that is sent when a connection is opened.
	CO MessageType = "CO"

	// ConnectionCloseMessage is a message that is sent when a connection is closed.
	CC MessageType = "CC"

	// ConnectionAcceptMessage is a message that is sent when a connection is accepted.
	CA MessageType = "CA"

	// ConnectionReadyMessage is a message that is sent after CA, when both sides of the connection are ready.
	CR MessageType = "CR"

	// ConnectionFailedMessage is a message that is sent when a connection open attempt failed.
	CF MessageType = "CF"

	// ConnectionNotFoundMessage is a message that is sent when a connection is not found.
	NF MessageType = "NF"

	// DataMessage is a message that contains data.
	D MessageType = "D"

	// DataGram a message that contains data from *gram socket.
	DG MessageType = "DG"

	// DataAckMessage is a message that is sent when data with a sequence number is received.
	DA MessageType = "DA"
)

// BridgeOptions defines the options for the bridge, which are shared with the relay on the other side of the bridge
// when this relay attempts to open a connection to the other relay.
type BridgeOptions struct {
	// Timestamp is the timestamp of the connection opening
	Timestamp time.Time

	// The remote URL
	URLRemote url.URL
}

type MessageHeader struct {
	// From is the spider device Id of the sender of the message
	From uuid.UUID

	// To is the spider device Id of the recipient of the message
	To uuid.UUID

	// The type of this message
	Type MessageType

	// CID is a uuid for the connection
	CID ConnectionID
}

// Message is a message that is sent to the portier server.
type Message struct {
	// Header is the plaintext, but authenticated header of the message
	Header MessageHeader

	// Message is the serialized and encrypted message, i.e. a DataMessage
	Message []byte
}

// ConnectionOpenMessage is a message that is sent when a connection is opened.
type ConnectionOpenMessage struct {
	// BridgeOptions defines the options for the bridge, which are shared
	BridgeOptions BridgeOptions
}

type ConnectionAcceptMessage struct {
}

// ConnectionFailedMessage is a message that is sent when a connection open attempt failed.
type ConnectionFailedMessage struct {
	// Reason is the reason why the connection failed
	Reason string
}

// DataMessage is a message that contains data.
type DataMessage struct {
	// Seq is the sequence number of the data
	Seq uint64

	// Retransmitted is a flag that indicates if the data is a retransmission
	Re bool

	// Data is the actual payload from the bridged connection
	Data []byte
}

// DataGramMessage is a message that contains data.
type DatagramMessage struct {
	// Addr is the address of the sender
	Source string

	// Target is the address of the recipient
	Target string

	// Data is the actual payload from the bridged connection
	Data []byte
}

// DataAckMessage is a message that is sent when data with a sequence number is received.
type DataAckMessage struct {
	// Seq is the sequence number of the data
	Seq uint64

	// Retransmitted is a flag that indicates if the ack is for a retransmitted message
	Re bool
}
