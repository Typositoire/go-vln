package proxy

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/typositoire/go-vln/backend"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	log "github.com/sirupsen/logrus"
)

// PClient ""
type PClient interface {
	Run()
}

type pClient struct {
	httpClient *resty.Client
	backend    backend.Backend
	logger     *log.Entry
}

// NewProxyClient ...
func NewProxyClient(options map[string]string) (PClient, error) {
	logger := log.WithFields(log.Fields{
		"component": "proxy.proxy_client",
	})

	log.SetFormatter(&log.JSONFormatter{})

	client := resty.New()

	client.
		SetHostURL(options["hostURL"]).
		SetLogger(logger)

	be, err := backend.NewBackend(options)

	if err != nil {
		return nil, err
	}

	err = be.Auth()

	if err != nil {
		return nil, err
	}

	return pClient{
		httpClient: client,
		backend:    be,
		logger:     logger,
	}, nil
}

func (pc pClient) Run() {
	e := echo.New()
	pc.setupRoutes(e)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Server.Addr = ":1323"

	e.Logger.Fatal(gracehttp.Serve(e.Server))
}

func (pc pClient) setupRoutes(e *echo.Echo) {
	e.GET("/*", pc.processRequest)
	e.DELETE("/*", pc.processRequest)
	e.POST("/*", pc.processRequest)
	e.PUT("/*", pc.processRequest)
	e.OPTIONS("/*", pc.processRequest)
	e.HEAD("/*", pc.processRequest)
	e.PATCH("/*", pc.processRequest)
}

func (pc pClient) processRequest(c echo.Context) error {
	var (
		body     string
		resp     *resty.Response
		err      error
		realPath string
	)

	if b, err := ioutil.ReadAll(c.Request().Body); err == nil {
		body = string(b)
	}

	headers := c.Request().Header
	method := c.Request().Method
	uri := c.Request().RequestURI
	backendCanProcess := canProcess(c.Request())
	backendIsInit, err := pc.backend.BackendIsInit()

	if err != nil {
		return err
	}

	if !backendCanProcess || !backendIsInit {
		pc.logger.Infoln("Not a get request, not handling.")
		resp, err = pc.passToVault(uri, body, method, headers)
	} else {
		pc.logger.Infoln("Looking for symlinks.")
		realPath, err = pc.backend.FindTarget(c.Request().RequestURI)
		if err != nil {
			return err
		}

		if realPath == "" {
			pc.logger.Infoln("Empty symlink location, assuming straight proxy.")
			resp, err = pc.passToVault(uri, body, method, headers)
		} else {
			resp, err = pc.passToVault(realPath, body, method, headers)
		}
	}

	if err != nil {
		return err
	}

	return c.String(resp.StatusCode(), string(resp.Body()))
}

func canProcess(r *http.Request) bool {
	if r.Method != "GET" {
		return false
	}

	if strings.HasPrefix(r.RequestURI, "/v1/sys") {
		return false
	}

	return true
}
