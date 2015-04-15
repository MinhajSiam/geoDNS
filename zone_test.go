package main

import (
	"github.com/miekg/dns"
	. "launchpad.net/gocheck"
)

func (s *ConfigSuite) TestExampleComZone(c *C) {
	ex := s.zones["test.example.com"]

	// test.example.com was loaded
	c.Assert(ex.Labels, NotNil)

	c.Check(ex.Logging.StatHat, Equals, true)
	c.Check(ex.Logging.StatHatAPI, Equals, "abc-test")

	c.Check(ex.Labels["weight"].MaxHosts, Equals, 1)

	// Make sure that the empty "no.bar" zone gets skipped and "bar" is used
	label, qtype := ex.findLabels("bar", []string{"no", "europe", "@"}, qTypes{dns.TypeA})
	c.Check(label.Records[dns.TypeA], HasLen, 1)
	c.Check(label.Records[dns.TypeA][0].RR.(*dns.A).A.String(), Equals, "192.168.1.2")
	c.Check(qtype, Equals, dns.TypeA)

	label, qtype = ex.findLabels("", []string{"@"}, qTypes{dns.TypeMX})
	Mxs := label.Records[dns.TypeMX]
	c.Check(Mxs, HasLen, 2)
	c.Check(Mxs[0].RR.(*dns.MX).Mx, Equals, "mx.example.net.")
	c.Check(Mxs[1].RR.(*dns.MX).Mx, Equals, "mx2.example.net.")

	label, qtype = ex.findLabels("", []string{"dk", "europe", "@"}, qTypes{dns.TypeMX})
	Mxs = label.Records[dns.TypeMX]
	c.Check(Mxs, HasLen, 1)
	c.Check(Mxs[0].RR.(*dns.MX).Mx, Equals, "mx-eu.example.net.")
	c.Check(qtype, Equals, dns.TypeMX)

	// look for multiple record types
	label, qtype = ex.findLabels("www", []string{"@"}, qTypes{dns.TypeCNAME, dns.TypeA})
	c.Check(label.Records[dns.TypeCNAME], HasLen, 1)
	c.Check(qtype, Equals, dns.TypeCNAME)

	label, qtype = ex.findLabels("", []string{"@"}, qTypes{dns.TypeNS})
	Ns := label.Records[dns.TypeNS]
	c.Check(Ns, HasLen, 2)
	c.Check(Ns[0].RR.(*dns.NS).Ns, Equals, "ns1.example.net.")
	c.Check(Ns[1].RR.(*dns.NS).Ns, Equals, "ns2.example.net.")

	label, qtype = ex.findLabels("", []string{"@"}, qTypes{dns.TypeSPF})
	Spf := label.Records[dns.TypeSPF]
	c.Check(Spf, HasLen, 1)
	c.Check(Spf[0].RR.(*dns.SPF).Txt[0], Equals, "v=spf1 ~all")

	label, qtype = ex.findLabels("foo", []string{"@"}, qTypes{dns.TypeTXT})
	Txt := label.Records[dns.TypeTXT]
	c.Check(Txt, HasLen, 1)
	c.Check(Txt[0].RR.(*dns.TXT).Txt[0], Equals, "this is foo")

	label, qtype = ex.findLabels("weight", []string{"@"}, qTypes{dns.TypeTXT})
	Txt = label.Records[dns.TypeTXT]
	c.Check(Txt, HasLen, 2)
	c.Check(Txt[0].RR.(*dns.TXT).Txt[0], Equals, "w1000")
	c.Check(Txt[1].RR.(*dns.TXT).Txt[0], Equals, "w1")
}

func (s *ConfigSuite) TestExampleOrgZone(c *C) {
	ex := s.zones["test.example.org"]

	// test.example.org was loaded
	c.Assert(ex.Labels, NotNil)

	label, qtype := ex.findLabels("sub", []string{"@"}, qTypes{dns.TypeNS})
	c.Assert(qtype, Equals, dns.TypeNS)

	Ns := label.Records[dns.TypeNS]
	c.Check(Ns, HasLen, 2)
	c.Check(Ns[0].RR.(*dns.NS).Ns, Equals, "ns1.example.com.")
	c.Check(Ns[1].RR.(*dns.NS).Ns, Equals, "ns2.example.com.")

}
