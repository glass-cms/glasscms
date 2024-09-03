//go:build go1.22

// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime"
)

// Defines values for ErrorCode.
const (
	ParameterInvalid      ErrorCode = "parameter_invalid"
	ParameterMissing      ErrorCode = "parameter_missing"
	ProcessingError       ErrorCode = "processing_error"
	ResourceAlreadyExists ErrorCode = "resource_already_exists"
	ResourceMissing       ErrorCode = "resource_missing"
)

// Defines values for ErrorType.
const (
	ApiError            ErrorType = "api_error"
	InvalidRequestError ErrorType = "invalid_request_error"
)

// Error Error is the response model when an API call is unsuccessful.
type Error struct {
	Code    ErrorCode              `json:"code"`
	Details map[string]interface{} `json:"details"`
	Message string                 `json:"message"`
	Type    ErrorType              `json:"type"`
}

// ErrorCode defines model for ErrorCode.
type ErrorCode string

// ErrorType defines model for ErrorType.
type ErrorType string

// Item Item represents an individual content item.
type Item struct {
	Content     string                 `json:"content"`
	CreateTime  time.Time              `json:"create_time"`
	DeleteTime  *time.Time             `json:"delete_time,omitempty"`
	DisplayName string                 `json:"display_name"`
	Metadata    map[string]interface{} `json:"metadata"`
	Name        string                 `json:"name"`
	Properties  map[string]interface{} `json:"properties"`
	UpdateTime  time.Time              `json:"update_time"`
}

// ItemCreate Resource create operation model.
type ItemCreate struct {
	Content     string                 `json:"content"`
	CreateTime  time.Time              `json:"create_time"`
	DisplayName string                 `json:"display_name"`
	Metadata    map[string]interface{} `json:"metadata"`
	Properties  map[string]interface{} `json:"properties"`
	UpdateTime  time.Time              `json:"update_time"`
}

// ItemKey defines model for ItemKey.
type ItemKey = string

// ItemsCreateJSONRequestBody defines body for ItemsCreate for application/json ContentType.
type ItemsCreateJSONRequestBody = ItemCreate

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /items)
	ItemsCreate(w http.ResponseWriter, r *http.Request)

	// (GET /items/{name})
	ItemsGet(w http.ResponseWriter, r *http.Request, name ItemKey)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ItemsCreate operation middleware
func (siw *ServerInterfaceWrapper) ItemsCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ItemsCreate(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ItemsGet operation middleware
func (siw *ServerInterfaceWrapper) ItemsGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name ItemKey

	err = runtime.BindStyledParameterWithOptions("simple", "name", r.PathValue("name"), &name, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ItemsGet(w, r, name)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/items", wrapper.ItemsCreate)
	m.HandleFunc("GET "+options.BaseURL+"/items/{name}", wrapper.ItemsGet)

	return m
}
