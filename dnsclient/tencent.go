package dnsclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

// TencentDNS  is a Client instance resolving using TencentDNS DNS-over-HTTPS service
type TencentDNS struct {
	BaseURL string
}

// Lookup performs a DNS lookup using TencentDNS
func (c *TencentDNS) Lookup(name string, rType uint16) Response {

	client := http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("accept", "application/dns-json")

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("type", strconv.Itoa(int(rType)))
	q.Add("cd", "false") // ignore DNSSEC
	req.URL.RawQuery = q.Encode()
	// fmt.Println(req.URL.String())

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("TENCENT DNS RESPONSE BODY:\n%s\n", body)

	dnsRequestResponse := requestResponse{}
	err = json.Unmarshal(body, &dnsRequestResponse)
	if err != nil {
		log.Fatal(err)
	}

	fout := Response{}

	if len(dnsRequestResponse.Answer) <= 0 {
		return fout
	}

	fout.TTL = dnsRequestResponse.Answer[0].TTL
	fout.Data = dnsRequestResponse.Answer[0].Data
	fout.Status = dns.RcodeToString[dnsRequestResponse.Status]

	return fout
}
