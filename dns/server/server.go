package server

import (
	"fmt"
	"log"

	"github.com/Abbas-gheydi/hanoproxy/dns/records"
	"github.com/miekg/dns"
)

type ServerConfs struct {
	EnableRecursiveDnsServer bool
	UpStreamDnsServer        string
	HADomain                 string
	ListenPort               string
	ListenIP                 string
	TTL                      string
	UpdateInterval           int
}

var (
	Configs = &ServerConfs{EnableRecursiveDnsServer: true,
		UpStreamDnsServer: "8.8.8.8",
		HADomain:          "abbas.local",
		ListenPort:        "1053",
		ListenIP:          "0.0.0.0",
		TTL:               "5",
	}
)

func Start() {

	dnsServer()
}

func handleHaRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	if r.Question[0].Qtype == 1 {
		requestedHost := r.Question[0].Name

		responseIp, _, responseErr := records.DnsRecords.Lookup(requestedHost)

		if responseErr != nil {
			log.Println(responseErr)
			dns.HandleFailed(w, r)
		} else {

			responseARecord := fmt.Sprintf("%v		%v	IN	A	%v", requestedHost, Configs.TTL, responseIp)
			aRecord, _ := dns.NewRR(responseARecord)

			m.Answer = []dns.RR{aRecord}

		}

	}

	err := w.WriteMsg(m)
	if err != nil {
		log.Println(err)
	}
}

func recursiveRequest(w dns.ResponseWriter, r *dns.Msg) {
	replayMsg := new(dns.Msg)
	replayMsg.SetReply(r)

	requestMsg := new(dns.Msg)
	requestMsg.Id = dns.Id()
	requestMsg.RecursionDesired = true
	requestMsg.Question = make([]dns.Question, 1)
	requestMsg.Question[0] = dns.Question{r.Question[0].Name, r.Question[0].Qtype, r.Question[0].Qclass}

	c := new(dns.Client)
	mssageFromUpstream, _, err := c.Exchange(requestMsg, Configs.UpStreamDnsServer+":53")
	if err == nil {

		replayMsg.Answer = mssageFromUpstream.Answer
		err = w.WriteMsg(replayMsg)
		if err != nil {
			log.Println(err)
		}

	} else {

		log.Println(Configs.UpStreamDnsServer + " not answer")

		dns.HandleFailed(w, r)

	}

}

func dnsServer() {
	dns.HandleFunc(Configs.HADomain, handleHaRequest)
	if Configs.EnableRecursiveDnsServer {
		dns.HandleFunc(".", recursiveRequest)
		log.Println("use", Configs.UpStreamDnsServer, "as upstream server")
	}
	log.Println("HA domain is", Configs.HADomain)
	log.Println("DNS Server Started at ", Configs.ListenIP, Configs.ListenPort)
	server := &dns.Server{Addr: Configs.ListenIP + ":" + Configs.ListenPort, Net: "udp"}
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
