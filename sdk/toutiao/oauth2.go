package toutiao

import (
	"encoding/json"

	"snk.git.node1/dan/go_request"
)

type Refresh_token_req struct {
	App_id        uint64 `json:"app_id"`
	Secret        string `json:"secret"`
	Grant_type    string `json:"grant_type"`
	Refresh_token string `json:"refresh_token"`
}

type Refresh_token_resp struct {
	Code       int                     `json:"code"`
	Message    string                  `json:"message"`
	Data       Refresh_token_data_resp `json:"data"`
	Request_id string                  `json:"request_id"`
}

type Refresh_token_data_resp struct {
	Access_token             string `json:"access_token"`
	Expires_in               uint64 `json:"expires_in"`
	Refresh_token            string `json:"refresh_token"`
	Refresh_token_expires_in uint64 `json:"refresh_token_expires_in"`
}

func Do_refresh_token(http_client go_request.Raw_request, raw_req *Refresh_token_req, timeout_s int) (int, *Refresh_token_resp, error) {
	if raw_req == nil {
		return 0, nil, nil
	}

	byte_req, err := json.Marshal(raw_req)
	if err != nil {
		return 0, nil, err
	}

	status_code, byte_resp, err := http_client.Post("https://ad.oceanengine.com/open_api/oauth2/refresh_token/").Set_json([]string{"Content-Type", "application/json;charset=UTF-8"}...).Send_timeout(byte_req, nil, timeout_s)
	if err != nil {
		return 0, nil, err

	}

	if status_code > 299 {
		return status_code, nil, nil
	}

	raw_resp := &Refresh_token_resp{}
	if err := json.Unmarshal(byte_resp, raw_resp); err != nil {
		return 0, nil, err
	}

	return status_code, raw_resp, nil
}
