package goStrongswanVici

import (
	"fmt"
)

type TerminateRequest struct {
	Child    string
	Ike      string
	Child_id string
	Ike_id   string
	Timeout  string
	Loglevel string
}

// To be simple, kill a client that is connecting to this server. A client is a sa.
//Terminates an SA while streaming control-log events.
func (c *Client) Terminate(r *TerminateRequest) (err error) {
	err = handlePanic(func() (err error) {
		msg, err := c.Request("terminate", map[string]interface{}{
			"child":    r.Child,
			"ike":      r.Ike,
			"child_id": r.Child_id,
			"ike_id":   r.Ike_id,
			"timeout":  r.Timeout,
			"loglevel": r.Loglevel,
		})
		if err != nil {
			return
		}
		if msg["success"] != "yes" {
			return fmt.Errorf("[Terminate] %s", msg["errmsg"])
		}
		return
	})
	return
}
