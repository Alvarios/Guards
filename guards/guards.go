package guards

import (
	"github.com/Alvarios/guards/config"
	"github.com/Alvarios/guards/log"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Event struct {
	Id         int
	StatusText string
	Message    string
}

type Guards struct {
	*zerolog.Logger
	Config config.LogConfig
}

/**
//Create a new instance of guards
*/
func NewLogger(config config.LogConfig) *Guards {
	file, err := os.Create(config.LogFile())
	if err != nil {
		//		t.Errorf("Failed to create file : %s", err.Error())
		return nil
	}

	level := zerolog.InfoLevel
	if config.IsDebug() {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(file).With().Timestamp().Logger()

	return &Guards{&logger, config}
}

var (
	invalidRequest = Event{Id: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest), Message: "Invalid request %s"}

	internalErrorRequest = Event{Id: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError), Message: "Status internal  : %s"}

	unauthorizedRequest = Event{Id: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized), Message: "Unauhthorized request : %s"}

	okRequest = Event{Id: http.StatusOK, StatusText: http.StatusText(http.StatusOK), Message: "ok request :  %s"}

	createdRequest = Event{Id: http.StatusCreated, StatusText: http.StatusText(http.StatusCreated), Message: "Ctreated request %s"}
)

// Invalid request
func (g *Guards) InvalidRequest(err error, message string) {
	g.Log().Str("service", g.Config.ServiceID()).
		Err(err).
		Int("id", invalidRequest.Id).
		Str("error_message", invalidRequest.StatusText).
		Msg(message)
}

//Internal server error
func (g *Guards) InternalErrorRequest(err error, message string) {
	g.Log().
		Str("service", g.Config.ServiceID()).
		Err(err).Int("id", internalErrorRequest.Id).
		Str("error_message", internalErrorRequest.StatusText).
		Msg(message)
}

//Unauthorized request
func (g *Guards) UnauthorizedRequest(err error, message string) {
	g.Log().
		Str("service", g.Config.ServiceID()).
		Err(err).
		Int("id", unauthorizedRequest.Id).
		Str("error_message", invalidRequest.StatusText).
		Msg(message)
}

func getIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(strings.Replace(ip, " ", "", -1))
			realIP = realIP
			ipAddress = ip
		}
	}
	return ipAddress
}

// Middleware
func (g *Guards) Middleware(next http.HandlerFunc) http.Handler {
	start := time.Now()
	le := &log.LogEntry{}
	fn := func(w http.ResponseWriter, r *http.Request) {
		le.ReceivedTime = start
		le.RequestMethod = r.Method
		le.RequestURL = r.URL.String()
		le.UserAgent = r.UserAgent()
		le.Referer = r.Referer()
		le.Proto = r.Proto
		le.RemoteIP = getIPAdress(r)
		next.ServeHTTP(w, r)
		le.Latency = time.Since(start)
		//le.Status = w.Header().Get("StatusCode")
		// status cide
		if le.Status == 0 {
			le.Status = http.StatusOK
		}
		g.Info().
			Str("service", g.Config.ServiceID()).
			Time("received_time", le.ReceivedTime).
			Str("method", le.RequestMethod).
			Str("url", le.RequestURL).
			Int64("header_size", le.RequestHeaderSize).
			Int64("body_size", le.RequestBodySize).
			Str("agent", le.UserAgent).
			Str("referer", le.Referer).
			Str("proto", le.Proto).
			Str("remote_ip", le.RemoteIP).
			Str("server_ip", le.ServerIP).
			Int("status", le.Status).
			Int64("resp_header_size", le.ResponseHeaderSize).
			Int64("resp_body_size", le.ResponseBodySize).
			Dur("latency", le.Latency).
			Msg("")

	}

	return http.HandlerFunc(fn) // wrapper
}
