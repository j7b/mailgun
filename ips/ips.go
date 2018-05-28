// Package ips implements IP management.
/*
This endpoint is accessed via the Caller type,
created with the API function.

Documentation for this endpoint is at
https://documentation.mailgun.com/en/latest/api-ips.html
*/
package ips

import (
	"github.com/j7b/mailgun/client"
)

type caller client.Caller

// Caller calls ip methods.
type Caller struct {
	caller
}

// API returns a Caller for c.
func API(c client.Caller) Caller {
	return Caller{caller: c}
}

// Info is information about an IP.
type Info struct {
	IP        string `json:"ip"`        // "ip": "192.161.0.1",
	Dedicated bool   `json:"dedicated"` // "dedicated": true,
	RDNS      string `json:"rdns"`      // "rdns": "luna.mailgun.net"
}

// Assigned returns a list of IPs allocated to the API domain.
func (c Caller) Assigned() ([]string, error) {
	var o struct {
		Items []string `json:"items"`
	}
	// client.Decode(res, &o, err)
	return o.Items, c.caller.Get(`/domains`, c.Domain(), `ips`).
		Decode(&o)
}

// Unassign removes a dedicated IP from the IP pool.
func (c Caller) Unassign(ip string) error {
	return c.Delete(`/domains`, c.Domain(), `ips`, ip).Err()
}

// Assign assigns a dedicated IP to the IP pool.
func (c Caller) Assign(ip string) error {
	return c.Post(`/domains`, c.Domain(), `ips`).
		SetForm("ip", ip).
		Err()
}

// Get retrieves Info for a particular IP.
func (c Caller) Get(ip string) (*Info, error) {
	var info *Info
	return info, c.caller.Get(`/ips`, ip).Decode(&info)
}

// List retrieves a list of IPs, if dedicated
// is true, only dedicated IPs.
func (c Caller) List(dedicated bool) ([]string, error) {
	req := c.caller.Get(`/ips`)
	if dedicated {
		req.SetQuery("dedicated", "true")
	}
	var o struct {
		Items []string `json:"items,omitempty"`
	}
	return o.Items, req.Decode(&o)
}
