package goStrongswanVici

import (
	"fmt"
)

type Pool struct {
	PoolMapping map[string]PoolMapping `json:"pools"`
}

type PoolMapping struct {
	Addrs string `json:"addrs"`
}

func (c *ClientConn) LoadPool(ph Pool) error {
	requestMap := &map[string]interface{}{}

	err := ConvertToGeneral(ph.PoolMapping, requestMap)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error creating request: %v", err)
	}

	msg, err := c.Request("load-pool", *requestMap)
	fmt.Println(msg)
	if msg["success"] != "yes" {
		return fmt.Errorf("unsuccessful LoadPool: %v", msg["success"])
	}

	return nil
}
