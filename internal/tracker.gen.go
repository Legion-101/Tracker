// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package tracker

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Error defines model for Error.
type Error struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

// Task defines model for Task.
type Task struct {
	Description *string `json:"description,omitempty"`
	IdTask      int     `json:"idTask"`
	NameTask    string  `json:"nameTask"`
}

// Tasks defines model for Tasks.
type Tasks = []Task

// PostTasksJSONRequestBody defines body for PostTasks for application/json ContentType.
type PostTasksJSONRequestBody = Task

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Метод создает задание с заданным id в образовательной программе
	// (DELETE /task/{idTask})
	DeleteTaskIdTask(ctx echo.Context, idTask int) error
	// Метод получает задание по id образовательной программы
	// (GET /task/{idTask})
	GetTaskIdTask(ctx echo.Context, idTask int) error
	// Метод получает все задания образовательной программы
	// (GET /tasks)
	GetTasks(ctx echo.Context) error
	// Метод создает задание с заданным id в образовательной программе
	// (POST /tasks)
	PostTasks(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// DeleteTaskIdTask converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteTaskIdTask(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "idTask" -------------
	var idTask int

	err = runtime.BindStyledParameterWithLocation("simple", false, "idTask", runtime.ParamLocationPath, ctx.Param("idTask"), &idTask)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter idTask: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteTaskIdTask(ctx, idTask)
	return err
}

// GetTaskIdTask converts echo context to params.
func (w *ServerInterfaceWrapper) GetTaskIdTask(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "idTask" -------------
	var idTask int

	err = runtime.BindStyledParameterWithLocation("simple", false, "idTask", runtime.ParamLocationPath, ctx.Param("idTask"), &idTask)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter idTask: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTaskIdTask(ctx, idTask)
	return err
}

// GetTasks converts echo context to params.
func (w *ServerInterfaceWrapper) GetTasks(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTasks(ctx)
	return err
}

// PostTasks converts echo context to params.
func (w *ServerInterfaceWrapper) PostTasks(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostTasks(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.DELETE(baseURL+"/task/:idTask", wrapper.DeleteTaskIdTask)
	router.GET(baseURL+"/task/:idTask", wrapper.GetTaskIdTask)
	router.GET(baseURL+"/tasks", wrapper.GetTasks)
	router.POST(baseURL+"/tasks", wrapper.PostTasks)

}