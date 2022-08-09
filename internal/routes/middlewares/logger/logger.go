package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Logger variables
//const (
//	TagPid               = "pid"
//	TagTime              = "time"
//	TagReferer           = "referer"
//	TagProtocol          = "protocol"
//	TagPort              = "port"
//	TagIP                = "ip"
//	TagIPs               = "ips"
//	TagHost              = "host"
//	TagMethod            = "method"
//	TagPath              = "path"
//	TagURL               = "url"
//	TagUA                = "ua"
//	TagLatency           = "latency"
//	TagStatus            = "status"
//	TagResBody           = "resBody"
//	TagReqHeaders        = "reqHeaders"
//	TagQueryStringParams = "queryParams"
//	TagBody              = "body"
//	TagBytesSent         = "bytesSent"
//	TagBytesReceived     = "bytesReceived"
//	TagRoute             = "route"
//	TagError             = "error"
//	// TagHeader DEPRECATED: Use TagReqHeader instead
//	TagHeader     = "header:"
//	TagReqHeader  = "reqHeader:"
//	TagRespHeader = "respHeader:"
//	TagLocals     = "locals:"
//	TagQuery      = "query:"
//	TagForm       = "form:"
//	TagCookie     = "cookie:"
//	TagBlack      = "black"
//	TagRed        = "red"
//	TagGreen      = "green"
//	TagYellow     = "yellow"
//	TagBlue       = "blue"
//	TagMagenta    = "magenta"
//	TagCyan       = "cyan"
//	TagWhite      = "white"
//	TagReset      = "reset"
//)

// Color values
const (
	//cBlack   = "\u001b[90m"
	cRed     = "\u001b[91m"
	cGreen   = "\u001b[92m"
	cYellow  = "\u001b[93m"
	cBlue    = "\u001b[94m"
	cMagenta = "\u001b[95m"
	cCyan    = "\u001b[96m"
	cWhite   = "\u001b[97m"
	cReset   = "\u001b[0m"
)

// New creates a new middleware handler
func New(logger *logrus.Entry, config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Get timezone location
	tz, err := time.LoadLocation(cfg.TimeZone)
	if err != nil || tz == nil {
		cfg.timeZoneLocation = time.Local
	} else {
		cfg.timeZoneLocation = tz
	}

	// Check if format contains latency
	cfg.enableLatency = strings.Contains(cfg.Format, "${latency}")

	// Create correct time-format
	var timestamp atomic.Value
	timestamp.Store(time.Now().In(cfg.timeZoneLocation).Format(cfg.TimeFormat))

	// Update date/time every 750 milliseconds in a separate go routine
	if strings.Contains(cfg.Format, "${time}") {
		go func() {
			for {
				time.Sleep(cfg.TimeInterval)
				timestamp.Store(time.Now().In(cfg.timeZoneLocation).Format(cfg.TimeFormat))
			}
		}()
	}

	// Set variables
	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Set error handler once
		once.Do(func() {
			errHandler = c.App().ErrorHandler
		})

		var start, stop time.Time

		// Set latency start time
		if cfg.enableLatency {
			start = time.Now()
		}

		// Handle request, store err for logging
		chainErr := c.Next()

		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		// Set latency stop time
		if cfg.enableLatency {
			stop = time.Now()
		}

		// Default output when no custom Format or io.Writer is given
		if cfg.enableColors && cfg.Format == ConfigDefault.Format {
			// Format error if exist
			// formatErr := ""
			// if chainErr != nil {
			// 	formatErr =  chainErr.Error()
			// }

			// { method: 'GET', path: /api/v2, status: 200, body: {anyany} }

			latency := stop.Sub(start).Round(time.Millisecond)

			params := c.Locals("params")

			var reqBodyBytes []byte
			if params != nil {
				reqBodyBytes, err = json.Marshal(params)
				if err != nil {
					reqBodyBytes = []byte("{}")
				} else {
					var mapParams map[string]interface{}
					err = json.Unmarshal(reqBodyBytes, &mapParams)
					if err == nil {
						hiddenFields := []string{"password", "otp_code", "email_code", "phone_code", "verification_code"}

						for _, field := range hiddenFields {
							if _, ok := mapParams[field]; ok {
								mapParams[field] = "[HIDDEN]"
							}
						}

						reqBodyBytes, err = json.Marshal(mapParams)
						if err != nil {
							reqBodyBytes = []byte("{}")
						}
					}
				}
			} else {
				reqBodyBytes = []byte("{}")
			}

			isResJson := true
			resBodyCompactedBuffer := new(bytes.Buffer)
			err = json.Compact(resBodyCompactedBuffer, c.Response().Body())
			if err != nil {
				isResJson = false
			}
			var resBodyBytes []byte
			if isResJson {
				resBodyBytes = resBodyCompactedBuffer.Bytes()
			} else {
				resBodyBytes = c.Response().Body()
			}

			resBodyBytes = bytes.TrimPrefix(resBodyBytes, []byte("\""))
			resBodyBytes = bytes.TrimSuffix(resBodyBytes, []byte("\""))

			if len(resBodyBytes) > 10_000 {
				resBodyBytes = []byte("[TRUNCATED]")
			}

			logStr := fmt.Sprintf(
				`{"method": %q, "path": %q, "status": %d, "ip": %q, "latency": %q, "payload": %s, "response": %s }`,
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				c.Locals("remote_ip").(net.IP).String(),
				latency,
				reqBodyBytes,
				resBodyBytes,
			)

			switch c.Response().StatusCode() {
			case 401, 422, 404, 405:
				logger.Warn(logStr)
			case 500:
				logger.Error(logStr)
			default:
				logger.Infof(logStr)
			}

			// End chain
			return nil
		}

		return nil
	}
}
