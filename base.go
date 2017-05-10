package dnspod

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//Time Overwrite the original time.Time layout for json.Unmarshal
type Time time.Time

const timeFormart = "2006-01-02 15:04:05"

//MarshalJSON json interface
func (p Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(p).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

//UnmarshalJSON json interface
func (p *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*p = Time(now)
	return
}

//String interface
func (p Time) String() string {
	return time.Time(p).Format(timeFormart)
}

//Yes converting "yes"/"no" to bool
type Yes bool

//MarshalJSON json interface
func (p Yes) MarshalJSON() ([]byte, error) {
	if p {
		return []byte(`"yes"`), nil
	}
	return []byte(`"no"`), nil
}

//UnmarshalJSON json interface
func (p *Yes) UnmarshalJSON(data []byte) (err error) {
	*p = string(data) == `"yes"`
	return
}

//String interface
func (p Yes) String() string {
	if p {
		return "yes"
	}
	return "no"
}

//Enable converting "enable"/"disable" to bool
type Enable bool

//MarshalJSON json interface
func (p Enable) MarshalJSON() ([]byte, error) {
	if p {
		return []byte(`"enable"`), nil
	}
	return []byte(`"disable"`), nil
}

//UnmarshalJSON json interface
func (p *Enable) UnmarshalJSON(data []byte) (err error) {
	*p = string(data) == `"enable"`
	return
}

//String interface
func (p Enable) String() string {
	if p {
		return "enable"
	}
	return "disable"
}

//Enabled converting "0"/"1" to bool
type Enabled bool

//MarshalJSON json interface
func (p Enabled) MarshalJSON() ([]byte, error) {
	if p {
		return []byte(`"1"`), nil
	}
	return []byte(`"0"`), nil
}

//UnmarshalJSON json interface
func (p *Enabled) UnmarshalJSON(data []byte) (err error) {
	*p = string(data) == `"1"`
	return
}

//String interface
func (p Enabled) String() string {
	if p {
		return "1"
	}
	return "0"
}

//----------------------------------------------------------------------------------
//# base struct

//Status of dnspod result
type Status struct {
	Code      int    `json:"code,string"`
	Message   string `json:"message"`
	CreatedAt Time   `json:"created_at"`
}

//----------------------------------------------------------------------------------
func simpleHTTP(client *http.Client, method, url string, params url.Values) (res []byte, err error) {
	body := bytes.NewBufferString(params.Encode())
	// Create request
	req, err := http.NewRequest(method, url, body)
	// Headers
	req.Header.Add("User-Agent", "Go DDNS/1.0.0 (bigemon@foxmail.com)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, nil
}
