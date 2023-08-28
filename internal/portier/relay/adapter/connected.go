package adapter

import (
	"fmt"
	"time"

	"github.com/marinator86/portier-cli/internal/portier/relay/encoder"
	"github.com/marinator86/portier-cli/internal/portier/relay/messages"
	"github.com/marinator86/portier-cli/internal/portier/relay/uplink"
)

type connectedState struct {
	// options are the connection adapter options
	options ConnectionAdapterOptions

	// encoderDecoder is the encoder/decoder for msgpack
	encoderDecoder encoder.EncoderDecoder

	// forwader is the forwarder
	forwarder Forwarder

	// uplink is the uplink
	uplink uplink.Uplink

	// CR ticker
	CRticker *time.Ticker

	// datachannel is the data channel
	datachannel chan messages.DataMessage

	// errorchannel is the error channel
	errorchannel chan error
}

func (c *connectedState) Start() error {
	// start reading from the connection
	dataChannel, errorChannel, err := c.forwarder.Start()
	if err != nil {
		return err
	}
	c.datachannel = dataChannel
	c.errorchannel = errorChannel

	// start CR ticker
	c.CRticker = time.NewTicker(1000 * time.Millisecond)
	msg := messages.Message{
		Header: messages.MessageHeader{
			From: c.options.LocalDeviceId,
			To:   c.options.PeerDeviceId,
			Type: messages.CR,
			CID:  c.options.ConnectionId,
		},
		Message: []byte{},
	}
	go func() {
		for range c.CRticker.C {
			err := c.uplink.Send(msg)
			if err != nil {
				fmt.Printf("error sending CR message: %s\n", err)
			}
		}
	}()

	return nil
}

func (c *connectedState) Stop() error {
	// send connection close message
	msg := messages.Message{
		Header: messages.MessageHeader{
			From: c.options.LocalDeviceId,
			To:   c.options.PeerDeviceId,
			Type: messages.CC,
			CID:  c.options.ConnectionId,
		},
		Message: []byte{},
	}
	_ = c.uplink.Send(msg)
	return c.forwarder.Close()
}

func (c *connectedState) HandleMessage(msg messages.Message) (ConnectionAdapterState, error) {
	// decrypt the data
	if msg.Header.Type == messages.D {
		c.CRticker.Stop()
		// TODO send acks
		return nil, nil
	} else if msg.Header.Type == messages.DA {
		c.CRticker.Stop()
		// encode the data
		ackMessage, err := c.encoderDecoder.DecodeDataAckMessage(msg.Message)
		if err != nil {
			return nil, err
		}
		c.forwarder.Ack(ackMessage.Seq)
		return nil, nil
	} else if msg.Header.Type == messages.CC {
		err := c.forwarder.Close()
		return nil, err
	} else if msg.Header.Type == messages.CR {
		return nil, nil
	}
	return nil, fmt.Errorf("expected message type [%s|%s|%s|%s], but got %s", messages.D, messages.DA, messages.CC, messages.CR, msg.Header.Type)
}

func NewConnectedState(options ConnectionAdapterOptions, uplink uplink.Uplink, forwarder Forwarder) ConnectionAdapterState {
	return &connectedState{
		options:        options,
		encoderDecoder: encoder.NewEncoderDecoder(),
		uplink:         uplink,
		forwarder:      forwarder,
	}
}