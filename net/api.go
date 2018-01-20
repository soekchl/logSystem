package net

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type Handler interface {
	HandleSession(*Session)
}

func Dial(network, address string, config *tls.Config, sendChanSize int) (*Session, error) {
	conn, err := tls.Dial(network, address, config)
	if err != nil {
		return nil, fmt.Errorf("[Dial] Error: %v", err)
	}

	return newSession(conn, sendChanSize), nil
}

func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if err != nil {
				err = fmt.Errorf("[Accept] Error: %v", err)
			}
			return nil, err
		}
		return conn, nil
	}
}
