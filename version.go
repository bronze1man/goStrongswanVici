package goStrongswanVici

type Version struct {
	Daemon  string
	Version string
	Sysname string
	Release string
	Machine string
}

func (c *Client) Version() (out *Version, err error) {
	err = handlePanic(func() (err error) {
		msg, err := c.Request("version", nil)
		if err != nil {
			return
		}
		out = &Version{
			Daemon:  msg["daemon"].(string),
			Version: msg["version"].(string),
			Sysname: msg["sysname"].(string),
			Release: msg["release"].(string),
			Machine: msg["machine"].(string),
		}
		return
	})
	return
}
