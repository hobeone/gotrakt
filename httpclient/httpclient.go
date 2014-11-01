/*
Package httpclient provides and easy to use http.Client that has support for
timeouts on the Dialing of the server as well as on reads and writes over
the connection.

Example usage:

Default settings:

client := httpclient.NewTimeoutClient()

Explicitly setting the connect and or the ReadWriteTimeout:

client := httpclient.NewTimeoutClient(
	httpclient.ConnectTimeout(2),
	httpclient.ReadWriteTimeout(5)
)

Use as normal:

client.Do(httpRequest)
*/

package httpclient

import (
	"net"
	"net/http"
	"time"
)

//TimeoutDialer is used bu the http.Client to Dial the server and sets the timeouts on reading and writing.
func (t *TimeoutClient) TimeoutDialer() func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, t.ConnectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(t.ReadWriteTimeout))
		return conn, nil
	}
}

//A TimeoutClient encapsulates handling timeouts for connecting and reading and writing with http.Client instances.
type TimeoutClient struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

type option func(*TimeoutClient)

// ConnectTimeout sets TimeoutClient's connect timeout to t seconds.
func ConnectTimeout(d time.Duration) option {
	return func(tc *TimeoutClient) {
		tc.ConnectTimeout = d
	}
}

// ReadWriteTimeout sets TimeoutClient's connect timeout to t seconds.
func ReadWriteTimeout(d time.Duration) option {
	return func(tc *TimeoutClient) {
		tc.ReadWriteTimeout = d
	}
}

// NewTimeoutClient returns a http.Client instance setup to use a TimeoutClient as the dial function
func NewTimeoutClient(options ...option) *http.Client {
	// Default configuration
	timeoutClient := &TimeoutClient{
		ConnectTimeout:   1 * time.Second,
		ReadWriteTimeout: 1 * time.Second,
	}
	for _, opt := range options {
		opt(timeoutClient)
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: timeoutClient.TimeoutDialer(),
		},
	}
}
