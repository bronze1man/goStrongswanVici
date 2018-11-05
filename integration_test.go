package goStrongswanVici

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	client, err := NewClientConnFromDefaultSocket()
	require.NoError(t, err)
	defer client.Close()

	// Get initial list of connections.
	initialConnections, err := client.ListConns("")
	require.NoError(t, err)

	// Create connection object.
	childConfMap := make(map[string]ChildSAConf)
	childSAConf := ChildSAConf{
		Local_ts:      []string{"10.10.59.0/24"},
		Remote_ts:     []string{"10.10.40.0/24"},
		ESPProposals:  []string{"aes256-sha256-modp2048"},
		StartAction:   "trap",
		CloseAction:   "restart",
		Mode:          "tunnel",
		ReqID:         "10",
		RekeyTime:     "10m",
		InstallPolicy: "no",
	}
	childConfMap["test-child-conn"] = childSAConf

	localAuthConf := AuthConf{
		AuthMethod: "psk",
	}
	remoteAuthConf := AuthConf{
		AuthMethod: "psk",
	}
	ikeConf := IKEConf{
		LocalAddrs:  []string{"192.168.198.10"},
		RemoteAddrs: []string{"192.168.198.11"},
		Proposals:   []string{"aes256-sha256-modp2048"},
		Version:     "1",
		LocalAuth:   localAuthConf,
		RemoteAuth:  remoteAuthConf,
		Children:    childConfMap,
		Encap:       "no",
	}
	ikeConfMap := map[string]IKEConf{"test-connection": ikeConf}

	// Add connection.
	err = client.LoadConn(&ikeConfMap)
	require.NoError(t, err)

	// Verify connection is added.
	connections, err := client.ListConns("")
	require.NoError(t, err)
	assert.Len(t, connections, len(initialConnections)+1)
}
