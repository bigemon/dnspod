package dnspod

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

//User information struct
type User struct {
	RealName          string  `json:"real_name"`
	UserType          string  `json:"user_type"`
	Telephone         string  `json:"telephone"`
	Nick              string  `json:"nick"`
	ID                int64   `json:"id,string"`
	Email             string  `json:"email"`
	Status            string  `json:"status"`
	EmailVerified     string  `json:"email_verified"`
	TelephoneVerified string  `json:"telephone_verified"`
	WeixinBinded      string  `json:"weixin_binded"`
	AgentPending      bool    `json:"agent_pending"`
	Balance           float32 `json:"balance"`
	Smsbalance        int     `json:"smsbalance"`
	UserGrade         string  `json:"user_grade"`
}

//UserAPI packaged some dnspod user APIs
type UserAPI struct {
	loginToken string
	client     *http.Client
}

//Detail Get Account Detail
func (p *UserAPI) Detail() (u User, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")
	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/User.Detail", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status `json:"status"`
		Info   struct {
			User User `json:"user"`
		} `json:"info"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return u, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Info.User, nil
}

//ModifyDetailOpt is opt arg of User.Modify
type ModifyDetailOpt struct {
	RealName  string //if no change,keep it null
	Nick      string //if no change,keep it null
	Telephone string //if no change,keep it null
}

//ModifyDetail use to modify detail of user
//opt:			The details you want to modify (RealName/Nick/Telephone)
func (p *UserAPI) ModifyDetail(opt ModifyDetailOpt) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	if opt.RealName != "" {
		params.Set("real_name", opt.RealName)
	}
	if opt.Nick != "" {
		params.Set("nick", opt.Nick)
	}
	if opt.Telephone != "" {
		params.Set("telephone", opt.Telephone)
	}

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/User.Modify", params)
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

//ModifyPassword use to modify password of user
//oldPwd:		Old password
//newPwd:		The new password you want to change
func (p *UserAPI) ModifyPassword(oldPwd, newPwd string) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("old_password", oldPwd)
	params.Set("new_password", newPwd)

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Userpasswd.Modify", params)
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

//ModifyEmail use to modify email of user
func (p *UserAPI) ModifyEmail(pwd, oldEmail, newEmail string) (err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("password", pwd)
	params.Set("old_email", oldEmail)
	params.Set("new_email", newEmail)

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Useremail.Modify", params)
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

//VerifyInfo Verification code and binding information
type VerifyInfo struct {
	Code string `json:"verify_code"`
	Desc string `json:"verify_desc"`
}

//PhoneVerify Get cell phone Verification code and binding information
func (p *UserAPI) PhoneVerify(phone string) (v VerifyInfo, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	params.Set("telephone", phone)
	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/Telephoneverify.Code", params)
	if err != nil {
		return
	}
	var jsonRes struct {
		Status Status     `json:"status"`
		Info   VerifyInfo `json:"user"`
	}
	if err = json.Unmarshal(res, &jsonRes); err != nil {
		return
	}
	if jsonRes.Status.Code != 1 {
		return v, errors.New(jsonRes.Status.Message)
	}
	return jsonRes.Info, nil
}

//Log get action logs
func (p *UserAPI) Log() (log []string, err error) {
	params := url.Values{}
	params.Set("login_token", p.loginToken)
	params.Set("format", "json")

	res, err := simpleHTTP(p.client, "POST", "https://dnsapi.cn/User.Log", params)
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
