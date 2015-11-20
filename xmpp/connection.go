// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Conn represents a connection to an XMPP server.
type Conn struct {
	config Config

	in           *xml.Decoder
	out          io.Writer
	rawOut       io.WriteCloser // doesn't log. Used for <auth>
	keepaliveOut io.Writer

	jid  string
	Rand io.Reader

	lock          sync.Mutex
	inflights     map[Cookie]inflight
	customStorage map[xml.Name]reflect.Type

	lastPingRequest  time.Time
	lastPongResponse time.Time

	delayedClose chan bool
	closed       bool
}

// NewConn creates a new connection
func NewConn(in *xml.Decoder, out io.Writer, jid string) *Conn {
	return &Conn{
		in:  in,
		out: out,
		jid: jid,

		inflights:    make(map[Cookie]inflight),
		delayedClose: make(chan bool),
	}
}

// Close closes the underlying connection
func (c *Conn) Close() error {
	if c.closed {
		return errors.New("xmpp: the connection is already closed")
	}

	c.closed = true

	// RFC 6120, Section 4.4 and 9.1.5
	log.Println("xmpp: sending closing stream tag")
	fmt.Fprint(c.out, "</stream:stream>")

	//TODO: find a better way to prevent sending message.
	c.out = ioutil.Discard

	select {
	case <-c.delayedClose:
	case <-time.After(30 * time.Second):
		log.Println("xmpp: timed out waiting for closing stream")
	}

	return c.closeTCP()
}

func (c *Conn) receivedClosingStreamTag() {
	go c.Close()
	c.delayedClose <- true
}

func (c *Conn) closeTCP() error {
	log.Println("xmpp: TCP closed")
	return c.rawOut.Close()
}

// Next reads stanzas from the server. If the stanza is a reply, it dispatches
// it to the correct channel and reads the next message. Otherwise it returns
// the stanza for processing.
func (c *Conn) Next() (stanza Stanza, err error) {
	for {
		if stanza.Name, stanza.Value, err = next(c); err != nil {
			return
		}

		if _, ok := stanza.Value.(*StreamClose); ok {
			log.Println("xmpp: received closing stream tag")
			go c.receivedClosingStreamTag()
			return
		}

		if iq, ok := stanza.Value.(*ClientIQ); ok && (iq.Type == "result" || iq.Type == "error") {
			var cookieValue uint64
			if cookieValue, err = strconv.ParseUint(iq.ID, 16, 64); err != nil {
				err = errors.New("xmpp: failed to parse id from iq: " + err.Error())
				return
			}
			cookie := Cookie(cookieValue)

			c.lock.Lock()
			inflight, ok := c.inflights[cookie]
			c.lock.Unlock()

			if !ok {
				continue
			}
			if len(inflight.to) > 0 {
				// The reply must come from the address to
				// which we sent the request.
				if inflight.to != iq.From {
					continue
				}
			} else {
				// If there was no destination on the request
				// then the matching is more complex because
				// servers differ in how they construct the
				// reply.
				if len(iq.From) > 0 && iq.From != c.jid && iq.From != RemoveResourceFromJid(c.jid) && iq.From != domainFromJid(c.jid) {
					continue
				}
			}

			c.lock.Lock()
			delete(c.inflights, cookie)
			c.lock.Unlock()

			inflight.replyChan <- stanza
			continue
		}

		return
	}
}

// Cancel cancels and outstanding request. The request's channel is closed.
func (c *Conn) Cancel(cookie Cookie) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	inflight, ok := c.inflights[cookie]
	if !ok {
		return false
	}

	delete(c.inflights, cookie)
	close(inflight.replyChan)
	return true
}

// Send sends an IM message to the given user.
func (c *Conn) Send(to, msg string) error {
	archive := ""
	if !c.config.Archive {
		// The first part of archive is from google:
		// See https://developers.google.com/talk/jep_extensions/otr
		// The second part of the stanza is from XEP-0136
		// http://xmpp.org/extensions/xep-0136.html#pref-syntax-item-otr
		// http://xmpp.org/extensions/xep-0136.html#otr-nego
		archive = "<nos:x xmlns:nos='google:nosave' value='enabled'/><arc:record xmlns:arc='http://jabber.org/protocol/archive' otr='require'/>"
	}
	_, err := fmt.Fprintf(c.out, "<message to='%s' from='%s' type='chat'><body>%s</body>%s</message>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(msg), archive)
	return err
}
