package goStrongswanVici

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type keyPayload struct {
	Typ  string `json:"type"`
	Data string `json:"data"`
}

func (c *ClientConn) LoadECDSAPrivateKey(key *ecdsa.PrivateKey) error {
	mk, err := x509.MarshalECPrivateKey(key)

	if err != nil {
		return err
	}

	var pemData = pem.EncodeToMemory(&pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: mk,
	})

	return c.LoadPrivateKey("ECDSA", string(pemData))
}

func (c *ClientConn) LoadRSAPrivateKey(key *rsa.PrivateKey) error {
	var mk = x509.MarshalPKCS1PrivateKey(key)

	var pemData = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: mk,
	})

	return c.LoadPrivateKey("RSA", string(pemData))
}

func (c *ClientConn) LoadPrivateKey(typ, data string) (err error) {
	requestMap := &map[string]interface{}{}

	var k = keyPayload{
		Typ:  typ,
		Data: data,
	}

	if err = ConvertToGeneral(k, requestMap); err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	msg, err := c.Request("load-key", *requestMap)
	if msg["success"] != "yes" {
		return fmt.Errorf("unsuccessful loadPrivateKey: %v", msg["success"])
	}

	return nil
}
