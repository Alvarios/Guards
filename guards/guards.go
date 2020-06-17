package guards

import (
	"github.com/Alvarios/guards/config"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Id         int
	StatusText string
	Message    string
}

type Guards struct {
	C alice.Chain
}

/**
//Create a new instance of guards
*/
func NewLogger(config config.LogConfig) *zerolog.Logger {
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

	logger := zerolog.
		New(file).
		With().
		//Str("role", "my-service").
		//Str("host", host).
		Timestamp().
		Logger()
	return &logger
}

func NewGuards(logger *zerolog.Logger) *Guards {
	c := alice.New()

	// Install the logger handler with default output on the console
	c = c.Append(hlog.NewHandler(*logger))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to that handler, all our logs will come with some prepopulated fields.
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	return &Guards{C: c}
}

var (
	invalidRequest = Event{Id: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest), Message: "Invalid request %s"}

	internalErrorRequest = Event{Id: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError), Message: "Status internal  : %s"}

	unauthorizedRequest = Event{Id: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized), Message: "Unauhthorized request : %s"}

	okRequest = Event{Id: http.StatusOK, StatusText: http.StatusText(http.StatusOK), Message: "ok request :  %s"}

	createdRequest = Event{Id: http.StatusCreated, StatusText: http.StatusText(http.StatusCreated), Message: "Ctreated request %s"}
)

// Invalid request
func (g *Guards) InvalidRequest(r *http.Request, err error, message string) {
	hlog.
		FromRequest(r).
		Error().
		//Str("service", g.Config.ServiceID()).
		Err(err).
		Int("id", invalidRequest.Id).
		Str("error_message", invalidRequest.StatusText).
		Msg(message)
}

/*
//Unauthorized request
func (g *Guards) UnauthorizedRequest(err error, message string) {
	g.Error().
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
		context := context.Background()
		ctx := g.With().Str("component", "module").Logger().WithContext(context)
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
}*/
