package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message, nil)
}

func InternalServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
}

// InvalidRequestMessage is a clean, client-safe substitute for the raw
// parser error gin's binding helpers (ShouldBindJSON, ShouldBindQuery,
// MultipartForm, ...) return on malformed or mistyped input - e.g. "json:
// cannot unmarshal number into Go struct field ..." - which otherwise leaks
// internal field/type details.
const InvalidRequestMessage = "invalid or malformed request data"

// SanitizedError logs the real error server-side (fully visible for
// debugging, tagged with a short context so it's traceable) and returns a
// generic, client-safe string to use in place of err.Error() in any 500
// response.
//
// This is only for genuinely unexpected/internal failures - a real
// validation or business-rule error (e.g. "room number already exists")
// is deliberately constructed by a service and already safe to show
// verbatim; those are unaffected and keep using err.Error() directly.
// What reaches a 500 branch in this codebase is, by construction, never
// one of those - it's always a lower-level failure (DB, driver, etc.)
// bubbling up from a pure read/passthrough call, so sanitizing every 500
// error message this way doesn't lose any real client-facing detail.
func SanitizedError(context string, err error) string {
	log.Printf("[internal error] %s: %v", context, err)
	return "internal server error"
}
