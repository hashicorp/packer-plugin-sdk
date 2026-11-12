// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/yamux"
)

// MuxBroker is responsible for brokering multiplexed connections by unique ID.
//
// This allows a plugin to request a channel with a specific ID to connect to
// or accept a connection from, and the broker handles the details of
// holding these channels open while they're being negotiated.
type muxBroker struct {
	nextId  uint32
	session *yamux.Session
	streams map[uint32]*muxBrokerPending

	sync.Mutex
}

type muxBrokerPending struct {
	ch       chan net.Conn
	doneCh   chan struct{}
	cancelCh chan struct{}

	cancelOnce sync.Once
	doneOnce   sync.Once
}

func (p *muxBrokerPending) cancel() {
	p.cancelOnce.Do(func() { close(p.cancelCh) })
}

func (p *muxBrokerPending) markDone() {
	p.doneOnce.Do(func() { close(p.doneCh) })
}

func newMuxBroker(s *yamux.Session) *muxBroker {
	return &muxBroker{
		session: s,
		streams: make(map[uint32]*muxBrokerPending),
	}
}

func newMuxBrokerClient(rwc io.ReadWriteCloser) (*muxBroker, error) {
	s, err := yamux.Client(rwc, nil)
	if err != nil {
		return nil, err
	}

	return newMuxBroker(s), nil
}

func newMuxBrokerServer(rwc io.ReadWriteCloser) (*muxBroker, error) {
	s, err := yamux.Server(rwc, nil)
	if err != nil {
		return nil, err
	}

	return newMuxBroker(s), nil
}

// Accept accepts a connection by ID.
//
// This should not be called multiple times with the same ID at one time.
func (m *muxBroker) Accept(id uint32) (net.Conn, error) {
	var c net.Conn
	p := m.getStream(id)
	select {
	case c = <-p.ch:
		p.markDone()
	case <-p.cancelCh:
		select {
		case conn := <-p.ch:
			conn.Close()
		default:
		}
		m.Lock()
		if cur, ok := m.streams[id]; ok && cur == p {
			delete(m.streams, id)
		}
		m.Unlock()
		return nil, fmt.Errorf("accept cancelled")
	case <-time.After(5 * time.Second):
		m.Lock()
		defer m.Unlock()
		delete(m.streams, id)

		return nil, fmt.Errorf("timeout waiting for accept")
	}

	// Ack our connection
	if err := binary.Write(c, binary.LittleEndian, id); err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

// Close closes the connection and all sub-connections.
func (m *muxBroker) Close() error {
	return m.session.Close()
}

// Dial opens a connection by ID.
func (m *muxBroker) Dial(id uint32) (net.Conn, error) {
	// Open the stream
	stream, err := m.session.OpenStream()
	if err != nil {
		return nil, err
	}

	// Write the stream ID onto the wire.
	if err := binary.Write(stream, binary.LittleEndian, id); err != nil {
		stream.Close()
		return nil, err
	}

	// Read the ack that we connected. Then we're off!
	var ack uint32
	if err := binary.Read(stream, binary.LittleEndian, &ack); err != nil {
		stream.Close()
		return nil, err
	}
	if ack != id {
		stream.Close()
		return nil, fmt.Errorf("bad ack: %d (expected %d)", ack, id)
	}

	return stream, nil
}

// CancelAccept stops a waiting Accept and closes a stream that already arrived. The shared session stays open.
func (m *muxBroker) CancelAccept(id uint32) {
	p := m.getStream(id)
	p.cancel()
	p.markDone()
	select {
	case conn := <-p.ch:
		conn.Close()
	default:
	}
	// Leave p in place so Accept sees the cancel instead of waiting on a new stream.
}

// NextId returns a unique ID to use next.
func (m *muxBroker) NextId() uint32 {
	return atomic.AddUint32(&m.nextId, 1)
}

// Run starts the brokering and should be executed in a goroutine, since it
// blocks forever, or until the session closes.
func (m *muxBroker) Run() {
	for {
		stream, err := m.session.AcceptStream()
		if err != nil {
			// Once we receive an error, just exit
			break
		}

		// Read the stream ID from the stream
		var id uint32
		if err := binary.Read(stream, binary.LittleEndian, &id); err != nil {
			stream.Close()
			continue
		}

		// Initialize the waiter
		p := m.getStream(id)
		select {
		case p.ch <- stream:
		default:
		}

		// Wait for a timeout
		go m.timeoutWait(id, p)
	}
}

func (m *muxBroker) getStream(id uint32) *muxBrokerPending {
	m.Lock()
	defer m.Unlock()

	p, ok := m.streams[id]
	if ok {
		return p
	}

	m.streams[id] = &muxBrokerPending{
		ch:       make(chan net.Conn, 1),
		doneCh:   make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	return m.streams[id]
}

func (m *muxBroker) timeoutWait(id uint32, p *muxBrokerPending) {
	// Wait for the stream to either be picked up and connected, or
	// for a timeout.
	timeout := false
	select {
	case <-p.doneCh:
	case <-time.After(5 * time.Second):
		timeout = true
	}

	// If we timed out, then check if we have a channel in the buffer,
	// and if so, close it. Non-blocking: Accept or CancelAccept may
	// already have taken the conn.
	if timeout {
		select {
		case s := <-p.ch:
			s.Close()
		default:
		}
	}

	// A cancelled accept must stay registered until Accept observes it.
	// Deleting it here would make that Accept wait on a new stream.
	select {
	case <-p.cancelCh:
		return
	default:
	}

	m.Lock()
	defer m.Unlock()
	if cur, ok := m.streams[id]; ok && cur == p {
		delete(m.streams, id)
	}
}
