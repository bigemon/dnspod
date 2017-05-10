package dnspod

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

//DomainAPI packaged some dnspod domain APIs
type DomainAPI struct {
	loginToken string
	client     *http.Client
}

//DomainInfo is the struct of domain statistic info
type DomainInfo struct {
	DomainTotal   int `json:"domain_total"`
	AllTotal      int `json:"all_total"`
	MineTotal     int `json:"mine_total"`
	ShareTotal    int `json:"share_total"`
	VIPTotal      int `json:"vip_total"`
	IsmarkTotal   int `json:"ismark_total"`
	PauseTotal    int `json:"pause_total"`
	ErrorTotal    int `json:"error_total"`
	LockTotal     int `json:"lock_total"`
	SPAMTotal     int `json:"spam_total"`
	VIPExpire     int `json:"vip_expire"`
	ShareOutTotal int `json:"share_out_total"`
}

//Domain Details
type Domain struct {
	ID               int64  `json:"id"`
	Status           Enable `json:"status"`
	Grade            string `json:"grade"`
	GroupID          int    `json:"group_id,string"`
	SearchEnginePush Yes    `json:"searchengine_push"`
	IsMark           Yes    `json:"is_mark"`
	TTL              int    `json:"ttl,string"`
	CnameSpeedup     Enable `json:"cname_speedup"`
	Remark           string `json:"remark"`
	CreatedOn        Time   `json:"created_on"`
	UpdatedOn        Time   `json:"updated_on"`
	Punycode         string `json:"punycode"`
	ExtStatus        string `json:"ext_status"`
	Name             string `json:"name"`
	GradeTitle       string `json:"grade_title"`
	IsVIP            Yes    `json:"is_vip"`
	Owner            string `json:"owner"`
	Records          int    `json:"records,string"`
	AuthToAnquanbao  bool   `json:"auth_to_anquanbao"`
}

//DomainCreateOpt is Optional arg of Domain.Create
type DomainCreateOpt struct {
	GroupID int
	IsMark  Yes
}

//Create a new domain record
func (p *DomainAPI) Create(domain string, opt ...DomainCreateOpt) (domainID int64, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	if len(opt) > 0 && opt[0].GroupID != 0 {
		params.Set("group_id", strconv.Itoa(opt[0].GroupID))
	}
	if len(opt) > 0 && opt[0].IsMark {
		params.Set("is_mark", opt[0].IsMark.String())
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Domain.Create", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
		Domain struct {
			ID int64 `json:"id,string"`
		} `json:"domain"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return 0, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Domain.ID, nil
}

//DLType is enum const of DomainListOpt.Type
type DLType string

const (
	//DLTypeAll list all domains
	DLTypeAll DLType = "all"
	//DLTypeMine list mine domain
	DLTypeMine DLType = "mine"
	//DLTypeShare list share domain
	DLTypeShare DLType = "share"
	//DLTypeIsMark list mark domain
	DLTypeIsMark DLType = "ismark"
	//DLTypePause  list pause domain
	DLTypePause DLType = "pause"
	//DLTypeVIP list VIP domain
	DLTypeVIP DLType = "vip"
	//DLTypeRecent list the most recently operated domain
	DLTypeRecent DLType = "recent"
	//DLTypeShareOut list the domain that are shared to others
	DLTypeShareOut DLType = "share_out"
)

//DomainListOpt is Optional arg of Domain.List
type DomainListOpt struct {
	Type    DLType `json:"type"`
	Offset  int    `json:"offset"`
	Length  int    `json:"length"`
	GroupID int    `json:"group_id"`
	Keyword string `json:"keyword"`
}

//List the domain records for the specified optional filter criteria
func (p *DomainAPI) List(opt ...DomainListOpt) (list []Domain, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")
	var o DomainListOpt
	if len(opt) > 0 {
		o = opt[0]
	}
	if o.GroupID != 0 {
		params.Set("group_id", strconv.Itoa(o.GroupID))
	}
	if o.Keyword != "" {
		params.Set("keyword", o.Keyword)
	}
	if o.Length != 0 {
		params.Set("length", strconv.Itoa(o.Length))
	}
	if o.Offset != 0 {
		params.Set("offset", strconv.Itoa(o.Offset))
	}
	if o.Type != "" {
		params.Set("type", string(o.Type))
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Domain.List", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status  Status   `json:"status"`
		Domains []Domain `json:"domains"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return list, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Domains, nil
}

//Remove a domain record
func (p *DomainAPI) Remove(domain string) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Domain.Remove", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return errors.New(jsonRes.Status.Message)
	}
	return nil
}

//Status setting a domain enabled or disable
func (p *DomainAPI) Status(domain string, enable Enable) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("status", enable.String())

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Domain.Status", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return errors.New(jsonRes.Status.Message)
	}
	return nil
}

//Info get a domain record
func (p *DomainAPI) Info(domain string) (info Domain, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Domain.Info", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
		Domain Domain `json:"domain"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return info, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Domain, nil
}

//DomainLogOpt is the optional arg of Domain.Log
type DomainLogOpt struct {
	Offset int
	Length int
}

//Log lists the log of the specified domain
func (p *DomainAPI) Log(domain string, opt ...DomainLogOpt) (log []string, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	var o DomainLogOpt
	if len(opt) > 0 {
		o = opt[0]
	}
	if o.Length != 0 {
		params.Set("length", strconv.Itoa(o.Length))
	}
	if o.Offset != 0 {
		params.Set("offset", strconv.Itoa(o.Offset))
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Domain.Log", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status   `json:"status"`
		Log    []string `json:"log"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return log, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Log, nil
}
