package goStrongswanVici

import (
	"fmt"
	"io"
	"net"
)

// This object is not thread safe.
// if you want concurrent, you need create more clients.
type Client struct {
	conn          net.Conn
	responseChan  chan segment
	eventHandlers map[string]func(response map[string]interface{})
	lastError     error
}

func (c *Client) Close() error {
	close(c.responseChan)
	c.lastError = io.ErrClosedPipe
	return c.conn.Close()
}

func NewClient(conn net.Conn) (client *Client) {
	client = &Client{
		conn:         conn,
		responseChan: make(chan segment, 2),
	}
	go client.readThread()
	return client
}

func NewClientFromDefaultSocket() (client *Client, err error) {
	conn, err := net.Dial("unix", "/var/run/charon.vici")
	if err != nil {
		return
	}
	return NewClient(conn), nil
}

func (c *Client) Request(apiname string, request map[string]interface{}) (response map[string]interface{}, err error) {
	err = writeSegment(c.conn, segment{
		typ:  stCMD_REQUEST,
		name: "version",
		msg:  request,
	})
	if err != nil {
		return
	}
	outMsg := <-c.responseChan
	if c.lastError != nil {
		return nil, c.lastError
	}
	if outMsg.typ != stCMD_RESPONSE {
		return nil, fmt.Errorf("[%s] response error %d", apiname, outMsg.typ)
	}
	return outMsg.msg, nil
}

func (c *Client) RegisterEvent(name string, handler func(response map[string]interface{})) (err error) {
	if c.eventHandlers[name] != nil {
		return fmt.Errorf("[event %s] register a event twice.", name)
	}
	c.eventHandlers[name] = handler
	err = writeSegment(c.conn, segment{
		typ:  stEVENT_REGISTER,
		name: name,
	})
	if err != nil {
		delete(c.eventHandlers, name)
		return
	}
	outMsg := <-c.responseChan
	if c.lastError != nil {
		delete(c.eventHandlers, name)
		return c.lastError
	}

	if outMsg.typ != stEVENT_CONFIRM {
		delete(c.eventHandlers, name)
		return fmt.Errorf("[event %s] response error %d", name, outMsg.typ)
	}
	return nil
}

func (c *Client) UnregisterEvent(name string) (err error) {
	err = writeSegment(c.conn, segment{
		typ:  stEVENT_UNREGISTER,
		name: name,
	})
	if err != nil {
		return
	}
	outMsg := <-c.responseChan
	if c.lastError != nil {
		return c.lastError
	}

	if outMsg.typ != stEVENT_CONFIRM {
		return fmt.Errorf("[event %s] response error %d", name, outMsg.typ)
	}
	delete(c.eventHandlers, name)
	return nil
}

func (c *Client) readThread() {
	for {
		outMsg, err := readSegment(c.conn)
		if err != nil {
			c.lastError = err
			return
		}
		switch outMsg.typ {
		case stCMD_RESPONSE, stEVENT_CONFIRM:
			c.responseChan <- outMsg
		case stEVENT:
			handler := c.eventHandlers[outMsg.name]
			if handler != nil {
				handler(outMsg.msg)
			}
		default:
			c.lastError = fmt.Errorf("[Client.readThread] unknow msg type %d", outMsg.typ)
			return
		}
	}
}
