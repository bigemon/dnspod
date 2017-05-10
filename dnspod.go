package dnspod

import (
	"net/http"
	"regexp"
)

//NewDnspod creates and initializes a new Dnspod instance
func NewDnspod(loginToken string) *Dnspod {
	hc := &http.Client{}
	R := RecordAPI{
		loginToken: loginToken,
		client:     hc,
	}
	U := UserAPI{
		loginToken: loginToken,
		client:     hc,
	}
	D := DomainAPI{
		loginToken: loginToken,
		client:     hc,
	}
	return &Dnspod{
		client: hc,
		Record: R,
		User:   U,
		Domain: D,
	}
}

//Dnspod api 1.0
type Dnspod struct {
	client *http.Client
	Record RecordAPI
	User   UserAPI
	Domain DomainAPI
}

//MyWANIP used to get your WAN IP
func (p *Dnspod) MyWANIP() (ip string, err error) {
	res, err := simpleHTTP(p.client, "GET", "http://m.tool.chinaz.com/ipsel", nil)
	if err != nil {
		return
	}
	return string(regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).Find(res)), nil
}
