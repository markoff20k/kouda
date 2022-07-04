package recover

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg"
	"github.com/zsmartex/pkg/log"

	"github.com/zsmartex/kouda/params"
	"github.com/zsmartex/kouda/utils"
)

func defaultStackTraceHandler(c *fiber.Ctx, e interface{}) {
	buf := make([]byte, 2048)
	buf = buf[:runtime.Stack(buf, false)]
	log.Errorf("Panic: %v\n%s\n", e, string(buf))
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				if cfg.EnableStackTrace {
					cfg.StackTraceHandler(c, r)
				}

				var ok bool
				if err, ok = r.(error); ok {
					var is_zsmart_pkg_error bool
					var is_fiber_pkg_error bool

					_, is_zsmart_pkg_error = err.(*pkg.Error)
					_, is_fiber_pkg_error = err.(*fiber.Error)

					if utils.IsNotFoundError(err) {
						err = params.ErrRecordNotFound
					} else if utils.IsDuplicateKeyError(err) {
						err_msg := err.Error()
						err_msg = utils.TrimStringBetween(err_msg, "index_", "(")
						err_msg = strings.TrimSuffix(err_msg, "\" ")

						columns_str := utils.TrimStringAfter(err_msg, "on_")

						columns := strings.Split(columns_str, "_and_")

						errors := make([]string, 0)
						for _, column := range columns {
							errors = append(errors, fmt.Sprintf("%s.taken", column))
						}

						err = pkg.NewError(fiber.StatusUnprocessableEntity, errors...)
					} else if !is_zsmart_pkg_error && !is_fiber_pkg_error {
						err = params.ErrServerInternal
					}
				} else {
					err = params.ErrServerInternal
				}
			}
		}()

		// Return err if exist, else move to next handler
		return c.Next()
	}
}
