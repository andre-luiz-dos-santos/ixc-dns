package main

import (
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/miekg/dns"
)

type dnsHandler struct {
	ixc        *IXC
	loginRE    *regexp.Regexp
	answerTTL  int
	filteredIP string
	notFoundIP string
	errorIP    string
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) <= 0 {
		return
	}

	msg := &dns.Msg{}
	msg.SetReply(r)

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true

		domain := r.Question[0].Name
		addr, err := h.getAddr(domain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "domain %v: %v\n", domain, err)
			// do not return
		}
		if addr == "" {
			if err != nil {
				msg.Rcode = dns.RcodeServerFailure
			} else {
				msg.Rcode = dns.RcodeNameError
			}
		} else {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   domain,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    uint32(h.answerTTL),
				},
				A: net.ParseIP(addr),
			})
		}
	}

	w.WriteMsg(msg)
}

func (h *dnsHandler) getAddr(domain string) (string, error) {
	match := h.loginRE.FindStringSubmatch(domain)
	if len(match) >= 2 {
		return h.getLoginAddr(match[1])
	}
	return h.filteredIP, nil
}

func (h *dnsHandler) getLoginAddr(login string) (string, error) {
	addr, err := h.ixc.GetAddrByLogin(login)
	if err != nil {
		return h.errorIP, fmt.Errorf("login %v: %w", login, err)
	}
	if addr == "" {
		// Login não existe ou não está conectado agora.
		return h.notFoundIP, nil
	}
	return addr, nil
}
