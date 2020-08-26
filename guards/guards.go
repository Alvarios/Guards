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
