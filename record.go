package dnspod

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

//---------------------------------------------------------------------------------

//RStatus is the emnum of the record state arg
type RStatus string

const (
	//RStatusEnable :Default value
	RStatusEnable RStatus = "enable"
	//RStatusDisable : When set to this value, parsing will not take effect or load balancing
	RStatusDisable RStatus = "disable"
)

//RType is the emnum of the record type arg
type RType string

const (
	//RTypeA is Address record, used to specify the IPV4 address of the domain name (such as: 8.8.8.8)
	RTypeA RType = "A"
	//RTypeCNAME use to point a domain name to another domain name, and then provide an IP address by another domain name, you need to add a CNAME record.
	RTypeCNAME RType = "CNAME"
	//RTypeMX If you need to set up mailboxes so that the mailboxes can receive messages, you need to add an MX record
	RTypeMX RType = "MX"
	//RTypeTXT Can fill in anything, length limit 255. Most TXT records are used to make SPF Records (anti-spam)
	RTypeTXT RType = "TXT"
	//RTypeNS Domain Name server records, if you need to handle domain name to other DNS service provider resolution, you need to add NS records
	RTypeNS RType = "NS"
	//RTypeAAAA A IPV6 address (for example, FF06:0:0:0:0:0:0:C3) that is used to specify a host name (or domain name) corresponding to the record.
	RTypeAAAA RType = "AAAA"
	//RTypeSRV Records which computer provides which service. The format is the name, point, and type of the protocol, such as _XMPP-SERVER._TCP.
	RTypeSRV RType = "SRV"
)

//RecordInfo info of record
type RecordInfo struct {
	SubDomains  int `json:"sub_domains,string"`
	RecordTotal int `json:"record_total,string"`
}

//RecordDomain domain info of record
type RecordDomain struct {
	ID        int64    `json:"id,string"`
	Name      string   `json:"name"`
	Punycode  string   `json:"punycode"`
	Grade     string   `json:"grade"`
	Owner     string   `json:"owner"`
	ExtStatus string   `json:"ext_status"`
	TTL       int      `json:"ttl"`
	MinTTL    int      `json:"min_ttl"`
	DnspodNS  []string `json:"dnspod_ns"`
	Status    string   `json:"status"`
}

//Record details
type Record struct {
	ID            int64   `json:"id,string"`
	TTL           int     `json:"ttl,string"`
	Value         string  `json:"value"`
	Enabled       Enabled `json:"enabled,string"`
	UpdatedOn     Time    `json:"updated_on"`
	Name          string  `json:"name"`
	Line          string  `json:"line"`
	LineID        string  `json:"line_id"`
	Type          string  `json:"type"`
	Weight        int     `json:"weight"`
	MonitorStatus string  `json:"monitor_status"`
	Remark        string  `json:"remark"`
	UseAQB        Yes     `json:"use_aqb"`
	MX            int     `json:"mx,string"`
}

// //RecordList all records under a domain
// type RecordList struct {
// 	Status  Status       `json:"status"`
// 	Domain  RecordDomain `json:"domain"`
// 	Info    RecordInfo   `json:"info"`
// 	Records []Record     `json:"records"`
// }

//----------------------------------------------------------------------------------

//RecordAPI packaged some dnspod record APIs
type RecordAPI struct {
	loginToken string
	client     *http.Client
}

//List is used to get a list of records for a specified domain
//domain: The domain name you want to get the dns record
func (p *RecordAPI) List(domain string) (list []Record, err error) {
	// Get DNS record list (POST https://dnsapi.cn/Record.List)
	params := url.Values{}
	params.Set("format", "json")
	params.Set("login_token", p.loginToken)
	params.Set("domain", domain)
	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.List", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status  Status   `json:"status"`
		Records []Record `json:"records"`
	}
	if err = json.Unmarshal(res, &list); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return list, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Records, nil
}

//DDNSOpt Opt arg struct
type DDNSOpt struct {
	SubDomain  string //The default value is "@"
	RecordLine string //The default value is "默认"
	Value      string //The default value is you wan ip
}

//DDNS used to update the specified DDNS record
//recordID: 	If you don't know the record ID, maybe you need to call RecordList first
//domain: 		The domain name you want to get the dns record
//opt:			The other optional arg
func (p *RecordAPI) DDNS(domain string, recordID int64, opt ...DDNSOpt) (err error) {
	var jsonRes struct {
		Status Status `json:"status"`
	}
	params := url.Values{}
	params.Set("record_id", strconv.FormatInt(recordID, 10))
	params.Set("domain", domain)
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")
	if len(opt) > 0 && opt[0].RecordLine != "" {
		params.Set("record_line", opt[0].RecordLine)
	} else {
		params.Set("record_line", "默认")
	}
	if len(opt) > 0 && opt[0].Value != "" {
		params.Set("value", opt[0].Value)
	}
	if len(opt) > 0 && opt[0].SubDomain != "" {
		params.Set("sub_domain", opt[0].SubDomain)
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Ddns", params)
	if err != nil {
		return
	}
	// fmt.Println("res:", string(res))
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return errors.New(jsonRes.Status.Message)
	}
	return nil
}

//RecordOpt Opt arg struct
type RecordOpt struct {
	SubDomain  string //The default value is "@"
	RecordLine string //The default value is "默认"
	Disable    bool   //If the incoming true, parsing does not take effect.
	MX         int    //MX priority, valid when the record type is MX, range 1-20.
	TTL        int    //Range 1-604800, different levels of domain names have different minimum values
	Weight     int    //Range 0-100, Available only in the Enterprise VIP domain, 0 is used to shut down
}

//Create used to create a DNS record
//domain: 		Domain name
//recordType:	Record type, uppercase , you can use RType series constants
//value: 		The value of the record(ip/mx/url...)
//opt:			The other optional arg
func (p *RecordAPI) Create(domain string, recordType RType, value string, opt ...RecordOpt) (id int64, err error) {
	if recordType == "MX" && (len(opt) == 0 || opt[0].MX == 0) {
		return 0, errors.New("Need to set up opt.MX")
	}
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("record_type", string(recordType))
	params.Set("value", value)
	var o RecordOpt
	if len(opt) > 0 {
		o = opt[0]
	}
	if o.RecordLine == "" {
		params.Set("record_line", "默认")
	} else {
		params.Set("record_line", o.RecordLine)
	}
	if o.Disable {
		params.Set("status", "disable")
	} else {
		params.Set("status", "enable")
	}
	if o.MX != 0 {
		params.Set("mx", strconv.Itoa(o.MX))
	}
	if o.TTL != 0 {
		params.Set("ttl", strconv.Itoa(o.TTL))
	}
	if o.Weight != 0 {
		params.Set("weight", strconv.Itoa(o.Weight))
	}
	if o.SubDomain != "" {
		params.Set("sub_domain", o.SubDomain)
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Create", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
		Record Record `json:"record"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return 0, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Record.ID, nil
}

//Modify used to modify a DNS record
//domain: 		Domain name
//recordID:		The specified record ID that you want to modify
//recordType:	Record type, uppercase
//value: 		The value of the record(ip/mx/url...)
//opt:			The other optional arg
func (p *RecordAPI) Modify(domain string, recordID int64, recordType RType, value string, opt ...RecordOpt) (err error) {
	if recordType == "MX" && (len(opt) == 0 || opt[0].MX == 0) {
		return errors.New("Need to set up opt.MX")
	}
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("record_id", strconv.FormatInt(recordID, 10))
	params.Set("record_type", string(recordType))
	params.Set("value", value)
	var o RecordOpt
	if len(opt) > 0 {
		o = opt[0]
	}
	if o.RecordLine == "" {
		params.Set("record_line", "默认")
	} else {
		params.Set("record_line", o.RecordLine)
	}
	if o.Disable {
		params.Set("status", "disable")
	} else {
		params.Set("status", "enable")
	}
	if o.MX != 0 {
		params.Set("mx", strconv.Itoa(o.MX))
	}
	if o.TTL != 0 {
		params.Set("ttl", strconv.Itoa(o.TTL))
	}
	if o.Weight != 0 {
		params.Set("weight", strconv.Itoa(o.Weight))
	}
	if o.SubDomain != "" {
		params.Set("sub_domain", o.SubDomain)
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Modify", params)
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

//Remove used to remove a DNS record
//domain: 		Domain name
//recordID:		The specified record ID that you want to remove
func (p *RecordAPI) Remove(domain string, recordID int64) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("record_id", strconv.FormatInt(recordID, 10))

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Remove", params)
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

//Remark used to remark a DNS record
//domain: 		Domain name
//recordID:		The specified record ID that you want to remove
//remark:		The remark you want to set,if set to null, will clear the remark
func (p *RecordAPI) Remark(domain string, recordID int64, remark string) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("record_id", strconv.FormatInt(recordID, 10))
	params.Set("remark", remark)

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Remark", params)
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

//Info used to get information from recordID
//domain: 		Domain name
//recordID:		The specified record ID that you want to get information
func (p *RecordAPI) Info(domain string, recordID int64) (r Record, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("record_id", strconv.FormatInt(recordID, 10))

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Info", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
		Record Record `json:"record"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return r, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Record, nil
}

//Status used to set new status of record
//domain: 		Domain name
//recordID:		The specified record ID that you want to get information
//enable:		Status of the record to set, if incoming false, parsing does not take effect.
func (p *RecordAPI) Status(domain string, recordID int64, enable Enable) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("domain", domain)
	params.Set("record_id", strconv.FormatInt(recordID, 10))
	params.Set("status", enable.String())

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Record.Status", params)
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
