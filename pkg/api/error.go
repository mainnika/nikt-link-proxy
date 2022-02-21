package api

import (
	"fmt"
	"net/http"

	routing "github.com/jackwhelpton/fasthttp-routing/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

//go:generate easyjson -all error.go

// Static implementation assertion:
var _ routing.HTTPError = (*HTTPError)(nil)
var _ error = (*HTTPError)(nil)

// HTTPError represents an error that happens during the api request
//easyjson:json
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewHTTPError constructs a new error and fills the message with http error if empty
func NewHTTPError(code int, message ...string) (e *HTTPError) {

	e = &HTTPError{Code: code}

	if len(message) > 0 {
		e.Message = message[0]
	} else {
		e.Message = e.Error()
	}

	return
}

// HTTPError implements error interface
func (v *HTTPError) Error() (message string) {

	message = v.Message
	if message == "" {
		message = fasthttp.StatusMessage(v.Code)
	}

	return
}

// StatusCode implements HTTPError interface and returns the HTTP status code.
func (v *HTTPError) StatusCode() (code int) {
	code = v.Code
	return
}

// ErrorNotFound renders http error-404 template
func (api *API) ErrorNotFound(c *routing.Context) (err error) {
	return NewHTTPError(http.StatusNotFound)
}

// UseErrorHandler is the middleware that catch handlers errors and render error template
func (api *API) UseErrorHandler(c *routing.Context) (err error) {

	worker := func() (err error) {

		// catch panic error
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			err = NewHTTPError(
				http.StatusInternalServerError,
				fmt.Sprintf("panic:\n%v", r),
			)
		}()

		err = c.Next()

		return
	}

	err = worker()
	if err == nil {
		return
	}

	c.Abort()

	logrus.Warnf("Cannot process request, %v", err)

	statusCode := http.StatusInternalServerError
	if httpError, isHttpError := err.(routing.HTTPError); isHttpError {
		statusCode = httpError.StatusCode()
	}

	c.SetStatusCode(statusCode)

	return c.Write(err)
}
