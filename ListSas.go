package goStrongswanVici

import (
	"fmt"
	"sync"
)

//from list-sa event
type IkeSa struct {
	Uniqueid        string //ike_id in terminate argument.
	Version         string
	State           string
	Local_host      string
	Local_id        string
	Remote_host     string
	Remote_id       string
	Remote_xauth_id string
	Initiator       string
	Initiator_spi   string
	Responder_spi   string
	Encr_alg        string
	Encr_keysize    string
	Integ_alg       string
	Integ_keysize   string
	Prf_alg         string
	Dh_group        string
	Established     string
	Rekey_time      string
	Reauth_time     string
	Child_sas       string
}

type Child_sas struct {
	Reqid string
	//... TODO finish it
}

// To be simple, list all clients that are connecting to this server .
// A client is a sa.
// Lists currently active IKE_SAs
func (c *Client) ListSas(ike string, ike_id string) (err error) {
	err = handlePanic(func() (err error) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		//register event
		err = c.RegisterEvent("list-sa", func(response map[string]interface{}) {
			fmt.Printf("%#v\n", response)
			wg.Done()
		})
		if err != nil {
			return
		}
		outMsg, err := c.Request("list-sas", map[string]interface{}{
			"ike":    ike,
			"ike_id": ike_id,
		})
		if err != nil {
			return
		}
		fmt.Printf("%#v\n", outMsg)
		wg.Wait()
		err = c.UnregisterEvent("list-sa")
		if err != nil {
			return
		}
		return
	})
	return
}
