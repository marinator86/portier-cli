package adapter

import (
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marinator86/portier-cli/internal/portier/relay/encoder"
	"github.com/marinator86/portier-cli/internal/portier/relay/messages"
	"github.com/marinator86/portier-cli/internal/portier/relay/uplink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConnection(testing *testing.T) {
	// GIVEN
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	// Signals
	connectionChannel := make(chan bool, 1)
	acceptedChannel := make(chan bool, 1)
	closedChannel := make(chan bool, 1)

	urlRemote, _ := url.Parse("tcp://localhost:" + fmt.Sprint(port))
	options := ConnectionAdapterOptions{
		ConnectionId:        "test-connection-id",
		LocalDeviceId:       uuid.New(),
		PeerDeviceId:        uuid.New(),
		PeerDevicePublicKey: "test-peer-device-public-key",
		BridgeOptions: messages.BridgeOptions{
			URLRemote: *urlRemote,
		},
	}
	encoderDecoder := encoder.NewEncoderDecoder()

	// mock uplink
	uplink := MockUplink{}

	uplink.On("Send", mock.MatchedBy(func(msg messages.Message) bool {
		if msg.Header.Type == messages.CA {
			acceptedChannel <- true
		}
		if msg.Header.Type == messages.CC {
			closedChannel <- true
		}
		return true
	})).Return(nil)

	underTest := NewConnectingInboundState(options, encoderDecoder, &uplink, 1000*time.Millisecond)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			testing.Errorf("expected err to be nil, got %v", err)
		}
		defer conn.Close()
		connectionChannel <- true
	}()

	// WHEN
	underTest.Start()

	// THEN
	<-connectionChannel // tcp connection established
	<-acceptedChannel   // connection accepted message sent
	<-acceptedChannel   // resend connection accepted message sent
	underTest.Stop()
	<-closedChannel // connection closed message sent
	uplink.AssertExpectations(testing)
}

func TestConnectionWithError(testing *testing.T) {
	// GIVEN
	port := 51222

	// Signals
	failedChannel := make(chan bool, 1)

	urlRemote, _ := url.Parse("tcp://localhost:" + fmt.Sprint(port))
	options := ConnectionAdapterOptions{
		ConnectionId:        "test-connection-id",
		LocalDeviceId:       uuid.New(),
		PeerDeviceId:        uuid.New(),
		PeerDevicePublicKey: "test-peer-device-public-key",
		BridgeOptions: messages.BridgeOptions{
			URLRemote: *urlRemote,
		},
	}
	encoderDecoder := encoder.NewEncoderDecoder()

	// mock uplink
	uplink := MockUplink{}

	uplink.On("Send", mock.MatchedBy(func(msg messages.Message) bool {
		if msg.Header.Type == messages.CF {
			// convert msg.Message to string
			msgText := string(msg.Message)
			assert.Contains(testing, msgText, fmt.Sprint(port))
			assert.Contains(testing, msgText, "connection refused")
			failedChannel <- true
		}
		return true
	})).Return(nil)

	underTest := NewConnectingInboundState(options, encoderDecoder, &uplink, 1000*time.Millisecond)

	// WHEN
	err := underTest.Start()

	// THEN
	assert.NotNil(testing, err)
	// assert error contains port and connection refused
	assert.Contains(testing, err.Error(), fmt.Sprint(port))
	assert.Contains(testing, err.Error(), "connection refused")
	<-failedChannel // connection failed message sent
	uplink.AssertExpectations(testing)
}

type MockUplink struct {
	mock.Mock
}

func (m *MockUplink) Connect() (chan []byte, error) {
	m.Called()
	return nil, nil
}

func (m *MockUplink) Send(message messages.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockUplink) Close() error {
	m.Called()
	return nil
}

func (m *MockUplink) Events() <-chan uplink.UplinkEvent {
	m.Called()
	return nil
}