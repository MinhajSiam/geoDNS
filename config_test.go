package main

import (
	"github.com/miekg/dns"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type ConfigSuite struct {
	zones Zones
}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestReadConfigs(c *C) {
	s.zones = make(Zones)
	configReadDir("dns", s.zones)

	// Just check that example.com loaded, too.
	c.Check(s.zones["example.com"].Origin, Equals, "example.com")

	tz := s.zones["test.example.com"]

	// The real tests are in test.example.com so we have a place
	// to make nutty configuration entries
>>>>>>>>>>>>>>>>>>>> File 1
	c.Check(s.zones["test.example.com"].Origin, Equals, "test.example.com")
	c.Check(s.zones["test.example.com"].Options.MaxHosts, Equals, 2)
        c.Check(s.zones["test.example.com"].Options.Contact, Equals, "support.bitnames.com")
	c.Check(s.zones["test.example.com"].Labels["weight"].MaxHosts, Equals, 1)
>>>>>>>>>>>>>>>>>>>> File 2
	c.Check(s.zones["test.example.com"].Origin, Equals, "test.example.com")
	c.Check(s.zones["test.example.com"].Options.MaxHosts, Equals, 2)
	c.Check(s.zones["test.example.com"].Labels["weight"].MaxHosts, Equals, 1)
>>>>>>>>>>>>>>>>>>>> File 3
	c.Check(tz.Origin, Equals, "test.example.com")
	c.Check(tz.Options.MaxHosts, Equals, 2)
	c.Check(tz.Options.Contact, Equals, "support.bitnames.com")
	c.Check(tz.Labels["weight"].MaxHosts, Equals, 1)

	/* test different cname targets */
	c.Check(tz.Labels["www"].
		firstRR(dns.TypeCNAME).(*dns.CNAME).
		Target, Equals, "geo.bitnames.com.")

	c.Check(tz.Labels["www-cname"].
		firstRR(dns.TypeCNAME).(*dns.CNAME).
		Target, Equals, "bar.test.example.com.")

	c.Check(tz.Labels["www-alias"].
		firstRR(dns.TypeCNAME).(*dns.CNAME).
		Target, Equals, "bar-alias.test.example.com.")

<<<<<<<<<<<<<<<<<<<<
}
