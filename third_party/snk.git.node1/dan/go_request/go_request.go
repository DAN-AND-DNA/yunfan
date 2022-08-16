package go_request

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type Config struct {
	Connect_timeout                  int
	Read_timeout                     int
	Write_timeout                    int
	Max_keepalive_idle_conn_duration int
	Max_keepalive_conn_duration      int // <= -1 unlimit
	Enable_keep_alive                bool
	Read_buffer_size                 int
	Dial                             func(addr string) (net.Conn, error)
}

type Raw_request struct {
	conf        Config
	http_client *fasthttp.Client

	url    string
	method string
	//headers  map[string]string
	headers  []string
	request  *fasthttp.Request
	response *fasthttp.Response
}

func New(config_list ...Config) Raw_request {
	var config Config
	len_config_list := len(config_list)
	if len_config_list >= 1 {
		config = config_list[len_config_list-1]
	} else {
		config = Config{
			Connect_timeout:                  3,
			Read_timeout:                     10,
			Write_timeout:                    10,
			Max_keepalive_idle_conn_duration: 10,
			Max_keepalive_conn_duration:      -1,
			Enable_keep_alive:                false,
			Dial:                             nil,
		}
	}

	if config.Connect_timeout <= 0 {
		config.Connect_timeout = 3
	}

	if config.Read_timeout <= 0 {
		config.Read_timeout = 0
	}

	if config.Write_timeout <= 0 {
		config.Write_timeout = 0
	}

	if config.Max_keepalive_idle_conn_duration <= 0 {
		config.Max_keepalive_idle_conn_duration = 0
	}

	if config.Read_buffer_size <= 0 {
		config.Read_buffer_size = 0
	}

	http_client := fasthttp.Client{
		Name:                      "37_men",
		ReadTimeout:               time.Duration(config.Read_timeout) * time.Second,
		WriteTimeout:              time.Duration(config.Write_timeout) * time.Second,
		MaxIdemponentCallAttempts: 0,
		MaxIdleConnDuration:       time.Duration(config.Max_keepalive_idle_conn_duration) * time.Second,
		ReadBufferSize:            config.Read_buffer_size,
		RetryIf: func(request *fasthttp.Request) bool {
			return false
		},

		Dial: func(addr string) (net.Conn, error) {
			if config.Dial == nil {
				conn, err := fasthttp.DialTimeout(addr, time.Duration(config.Connect_timeout)*time.Second)

				if err != nil {
					return nil, err
				}
				return conn, nil
			} else {
				return config.Dial(addr)
			}
		},
	}

	if config.Max_keepalive_conn_duration >= 0 {
		http_client.MaxConnDuration = time.Duration(config.Max_keepalive_conn_duration) * time.Second
	}

	return Raw_request{
		conf:        config,
		http_client: &http_client,
	}
}

func (obj Raw_request) Acquire() Raw_request {
	obj.request = fasthttp.AcquireRequest()
	obj.response = fasthttp.AcquireResponse()

	return obj
}

func (obj Raw_request) Get(url ...string) Raw_request {
	if obj.request == nil {
		obj.request = fasthttp.AcquireRequest()
	}

	if obj.response == nil {
		obj.response = fasthttp.AcquireResponse()
	}

	if len(url) != 0 {
		obj.url = url[0]
	}

	obj.method = "GET"

	return obj
}

func (obj Raw_request) Post(url ...string) Raw_request {
	if obj.request == nil {
		obj.request = fasthttp.AcquireRequest()
	}
	if obj.response == nil {
		obj.response = fasthttp.AcquireResponse()
	}

	if len(url) != 0 {
		obj.url = url[0]
	}

	obj.method = "POST"
	return obj
}

func (obj Raw_request) Put(url ...string) Raw_request {
	if obj.request == nil {
		obj.request = fasthttp.AcquireRequest()
	}
	if obj.response == nil {
		obj.response = fasthttp.AcquireResponse()
	}

	if len(url) != 0 {
		obj.url = url[0]
	}

	obj.method = "PUT"
	return obj
}

func (obj Raw_request) Set(headers ...string) Raw_request {
	len_headers := len(headers)
	if len_headers%2 != 0 {
		obj.headers = headers[0 : len_headers-2]
		len_headers -= 1
	}

	obj.headers = headers

	return obj
}

func (obj Raw_request) Set_json(headers ...string) Raw_request {
	headers = append(headers, fasthttp.HeaderContentType)
	headers = append(headers, "application/json")
	return obj.Set(headers...)
}

// FIXME the body is invalid after calling End()
func (obj Raw_request) Send_unsafe(body []byte, headers map[string]string, wait_s int, change_url ...string) (Raw_request, int, []byte, error) {
	obj.request.Reset()
	obj.response.Reset()

	request := obj.request
	response := obj.response
	len_headers := len(obj.headers)

	for i := 0; i < len_headers; i += 2 {
		request.Header.Set(obj.headers[i], obj.headers[i+1])
	}

	if len(change_url) != 0 {
		obj.url = change_url[0]
	}

	request.Header.SetRequestURI(obj.url)
	request.Header.SetMethod(obj.method)
	request.AppendBody(body)

	if wait_s <= 0 {
		err := obj.http_client.Do(request, response)
		if err != nil {
			return obj, 0, nil, err
		}
	} else {
		err := obj.http_client.DoTimeout(request, response, time.Duration(wait_s)*time.Second)
		if err != nil {
			return obj, 0, nil, err
		}
	}

	if headers != nil {
		response.Header.VisitAll(func(key, value []byte) {
			(headers)[string(key)] = string(value)
		})
	}

	return obj, response.StatusCode(), response.Body(), nil
}

// FIXME the body is valid after calling End()
func (obj Raw_request) Send(body []byte, headers map[string]string) (int, []byte, error) {
	new_obj, status_code, new_body, err := obj.Send_unsafe(body, headers, -1)
	defer new_obj.Release()

	if err != nil {
		return 0, nil, err
	} else {
		if status_code == fasthttp.StatusOK {
			len_new_body := len(new_body)
			if len_new_body > 0 {
				resp_body := make([]byte, len(new_body))
				copy(resp_body, new_body)
				return status_code, resp_body, nil
			}
		}
	}
	return status_code, nil, nil
}

func (obj Raw_request) Send_timeout(body []byte, headers map[string]string, wait_s int) (int, []byte, error) {
	new_obj, status_code, new_body, err := obj.Send_unsafe(body, headers, wait_s)
	defer new_obj.Release()

	if err != nil {
		return 0, nil, err
	} else {
		if status_code == fasthttp.StatusOK {
			len_new_body := len(new_body)
			if len_new_body > 0 {
				resp_body := make([]byte, len(new_body))
				copy(resp_body, new_body)
				return status_code, resp_body, nil
			}
		}
	}
	return status_code, nil, nil
}

func (obj Raw_request) Reuse() {
	obj.request.Reset()
	obj.response.Reset()
}

func (obj Raw_request) Release() {
	fasthttp.ReleaseRequest(obj.request)
	fasthttp.ReleaseResponse(obj.response)
}

//TODO stream
