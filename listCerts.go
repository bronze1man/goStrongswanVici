package goStrongswanVici

import (
	"fmt"
	"strconv"
)

type Cert struct {
	Type       string
	Flag       string
	HasPrivKey bool
	Data       string
	Subject    string
	NotBefore  string
	NotAfter   string
}

// typ = certificate type to filter for, X509|X509_AC|X509_CRL|OCSP_RESPONSE|PUBKEY  or ANY
// flag = X.509 certificate flag to filter for, NONE|CA|AA|OCSP or ANY
// subject = set to list only certificates having subject
func (c *ClientConn) ListCerts(typ, flag, subject string) ([]*Cert, error) {
	certs := []*Cert{}
	var eventErr error
	err := c.RegisterEvent("list-cert", func(response map[string]interface{}) {
		cert := &Cert{}
		cert.Type, _ = response["type"].(string)
		cert.Flag, _ = response["flag"].(string)
		cert.HasPrivKey, _ = response["hasPrivKey"].(bool)
		cert.Data, _ = response["data"].(string)
		cert.Subject, _ = response["subject"].(string)
		cert.NotBefore, _ = response["notBefore"].(string)
		cert.NotAfter, _ = response["notAfter"].(string)
		if hasPri, ok := response["hasPriKey"].(string); ok {
			cert.HasPrivKey, _ = strconv.ParseBool(hasPri)
		}
		certs = append(certs, cert)
	})
	if err != nil {
		return nil, fmt.Errorf("error registering list-cert event: %v", err)
	}
	if eventErr != nil {
		return nil, eventErr
	}
	reqMap := map[string]interface{}{
		"type":    typ,
		"flag":    flag,
		"subject": subject,
	}
	_, err = c.Request("list-certs", reqMap)
	if err != nil {
		return nil, fmt.Errorf("error requesting list-certs: %v", err)
	}

	err = c.UnregisterEvent("list-cert")
	if err != nil {
		return nil, fmt.Errorf("error unregistering list-cert event: %v", err)
	}
	return certs, nil
}
