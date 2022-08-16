package http

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"errors"
	"sync"
	"time"
)

type mode int

const (
	Mode_none       mode = 0
	Mode_deployment mode = 1
	Mode_production mode = 2
	Mode_test       mode = 3
)

var (
	default_server = New_server()

	ERR_IS_RUNNING      = errors.New("is already running")
	ERR_LISTEN_FAILED   = errors.New("listen failed: client try 3 times !!!")
	ERR_RUN_AS_TEST     = errors.New("should run as test")
	ERR_BAD_RUN_AS_MODE = errors.New("bad run as mode")
)

type Http_server struct {
	is_running    bool
	is_test       bool
	get_handlers  map[string]func(*fiber.Ctx) error
	post_handlers map[string]func(*fiber.Ctx) error
	app           *fiber.App
	m_server      sync.Mutex
	err_server    error
}

func New_server() *Http_server {
	return &Http_server{
		get_handlers:  make(map[string]func(*fiber.Ctx) error),
		post_handlers: make(map[string]func(*fiber.Ctx) error),
	}
}

func new_app(server_name string, body_limit_bytes, read_timeout_secs, write_timeout_secs, max_conns int, get_only bool) *fiber.App {

	if body_limit_bytes < 0 {
		body_limit_bytes = 1 * 1024 * 1024
	}

	if read_timeout_secs < 0 {
		read_timeout_secs = 10
	}

	if write_timeout_secs < 0 {
		write_timeout_secs = 10
	}

	if max_conns < 0 {
		max_conns = 1024
	}

	return fiber.New(fiber.Config{
		Prefork:          false,
		ServerHeader:     server_name,
		CaseSensitive:    true,
		BodyLimit:        body_limit_bytes,
		Concurrency:      max_conns,
		ReadTimeout:      time.Duration(read_timeout_secs) * time.Second,
		WriteTimeout:     time.Duration(write_timeout_secs) * time.Second,
		GETOnly:          get_only,
		DisableKeepalive: false,
	})
}

func (this *Http_server) get_error() error {
	this.m_server.Lock()
	defer this.m_server.Unlock()

	if this.err_server != nil {
		return errors.New(this.err_server.Error())
	}

	return nil
}

type Http_server_config struct {
	Run_as             string
	Port               string
	Allow_origins      string
	Server_name        string
	Body_limit_bytes   int
	Read_timeout_secs  int
	Write_timeout_secs int
	Max_conns          int
	Get_only           bool
}

func (this *Http_server) listen(config Http_server_config) error {

	if this.is_running {
		return ERR_IS_RUNNING
	}

	if config.Run_as == "" || config.Port == "" {
		return nil
	}

	this.app = new_app(config.Server_name, config.Body_limit_bytes, config.Read_timeout_secs, config.Write_timeout_secs, config.Max_conns, config.Get_only)

	run_as := Mode_none

	switch config.Run_as {
	case "deployment":
		run_as = Mode_deployment
	case "production":
		run_as = Mode_production
	default:
		return ERR_BAD_RUN_AS_MODE
	}

	if run_as == Mode_deployment {
		this.app.Use(logger.New(logger.Config{
			Next:       nil,
			Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: "2006/01/02 15:04:05",
			TimeZone:   "Local",
			Output:     os.Stdout,
		}))
	}

	this.app.Use(cors.New(cors.Config{
		AllowOrigins: config.Allow_origins,
		AllowMethods: "GET,POST,HEAD",
	}))

	str_test_path := "/are_u_ok"

	this.app.Get(str_test_path, func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	for path, h := range this.get_handlers {
		this.app.Get(path, h)
	}

	for path, h := range this.post_handlers {
		this.app.Post(path, h)
	}

	this.app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	this.is_running = false
	if run_as != Mode_test {
		this.is_test = false

		go func() {
			if err := this.app.Listen("0.0.0.0:" + config.Port); err != nil {
				this.m_server.Lock()
				this.err_server = err
				this.m_server.Unlock()
			}
		}()

		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			a := fiber.AcquireAgent()
			a = a.Timeout(1 * time.Second)
			req := a.Request()
			req.Header.SetMethod(fiber.MethodGet)
			req.SetRequestURI("http://127.0.0.1:" + config.Port + str_test_path)

			if err := a.Parse(); err != nil {
				return err
			}

			if _, _, err := a.Bytes(); err == nil {
				this.is_running = true
				break
			}
		}
	} else {
		this.is_running = true
		this.is_test = true
		this.app.Handler()
	}

	if !this.is_running {
		if err := this.get_error(); err != nil {
			return err
		} else {
			return ERR_LISTEN_FAILED
		}
	}

	fmt.Println("http server is running")
	return nil
}

func (this *Http_server) shutdown(deadline_secs int) {
	if !this.is_running || this.is_test {
		return
	}

	go func() {
		this.app.Shutdown()
	}()

	time.Sleep(time.Duration(deadline_secs) * time.Second)
}

func (this *Http_server) register(method, path string, ep endpoint.Endpoint, decode_req func(*fiber.Ctx) (interface{}, error), decode_resp func(*fiber.Ctx, interface{}) error) {

	if this.is_running || path == "" {
		return
	}

	switch strings.ToLower(method) {
	case "get":
		this.get_handlers[path] = make_http_handler(ep, decode_req, decode_resp)
	case "post":
		this.post_handlers[path] = make_http_handler(ep, decode_req, decode_resp)
	}
}

func make_http_handler(ep endpoint.Endpoint, decode_req func(*fiber.Ctx) (interface{}, error), decode_resp func(*fiber.Ctx, interface{}) error) func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		if req, err := decode_req(c); err != nil {
			if resp, err := ep(ctx, err); err != nil {
				return decode_resp(c, err)
			} else {
				return decode_resp(c, resp)
			}
		} else {
			if resp, err := ep(ctx, req); err != nil {
				return decode_resp(c, err)
			} else {
				return decode_resp(c, resp)
			}
		}
	}
}

func Listen(config Http_server_config) error {
	return default_server.listen(config)
}

func Shutdown(deadline_secs int) {
	default_server.shutdown(deadline_secs)
}

func Register(method, path string, ep endpoint.Endpoint, decode_req func(*fiber.Ctx) (interface{}, error), decode_resp func(*fiber.Ctx, interface{}) error) {
	default_server.register(method, path, ep, decode_req, decode_resp)
}

func Get_error() error {
	return default_server.get_error()
}
