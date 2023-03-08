package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"golang-standards-project-example/internal/pkg/server"
	"net"
	"strconv"
)

type HttpServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewHttpServingOptions() *HttpServingOptions {
	return &HttpServingOptions{
		BindAddress: "127.0.0.1",
		BindPort:    8080,
	}
}

func (h *HttpServingOptions) ApplyTo(c *server.Config) error {
	c.HttpServing = &server.HttpServingInfo{Address: net.JoinHostPort(h.BindAddress, strconv.Itoa(h.BindPort))}
	return nil
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (h *HttpServingOptions) Validate() []error {
	var errors []error

	if h.BindPort < 0 || h.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--http.bind-port %v must be between 0 and 65535, inclusive. 0 for turning off insecure (HTTP) port",
				h.BindPort,
			),
		)
	}

	return errors
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (h *HttpServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&h.BindAddress, "http.bind-address", h.BindAddress, ""+
		"The IP address on which to serve the --http.bind-port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	fs.IntVar(&h.BindPort, "http.bind-port", h.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")
}
