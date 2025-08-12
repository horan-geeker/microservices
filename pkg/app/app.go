package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"microservices/entity"
	errors2 "microservices/pkg/ecode"
	"microservices/pkg/log"
	"microservices/pkg/util"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var app *App

// App 包了一层 gin 内核
type App struct {
	command *command
	*gin.Engine
	serverOptions *ServerOptions
}

// POST is a shortcut for router.Handle("POST", path, handlers).
func (a *App) POST(relativePath string, handlers ...any) gin.IRoutes {
	return a.Handle(http.MethodPost, relativePath, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handlers).
func (a *App) GET(relativePath string, handlers ...any) gin.IRoutes {
	return a.Handle(http.MethodGet, relativePath, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handlers).
func (a *App) DELETE(relativePath string, handlers ...any) gin.IRoutes {
	return a.Handle(http.MethodDelete, relativePath, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handlers).
func (a *App) PATCH(relativePath string, handlers ...any) gin.IRoutes {
	return a.Handle(http.MethodPatch, relativePath, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handlers).
func (a *App) PUT(relativePath string, handlers ...any) gin.IRoutes {
	return a.Handle(http.MethodPut, relativePath, handlers...)
}

func (a *App) Handle(httpMethod, relativePath string, customHandlers ...any) gin.IRoutes {
	handlers := make([]gin.HandlerFunc, 0)
	for _, handler := range customHandlers {
		var ginHandlerFunc gin.HandlerFunc
		name := util.GetFunctionName(handler)
		if strings.Contains(name, "/controller/") {
			if err := a.validateController(handler); err != nil {
				panic(err)
			}
			ginHandlerFunc = a.wrapperGin(handler)
		} else {
			var ok bool
			ginHandlerFunc, ok = handler.(gin.HandlerFunc)
			if !ok {
				panic(ok)
			}
		}
		handlers = append(handlers, ginHandlerFunc)
	}
	return a.Engine.Handle(httpMethod, relativePath, handlers...)
}

func (a *App) wrapperGin(handle any) gin.HandlerFunc {
	t := reflect.TypeOf(handle)
	f := reflect.ValueOf(handle)
	argsNum := t.NumIn()
	var controllerParamsType []reflect.Kind
	var bodyStruct any
	// 控制器最后一个参数是结构体解析的 body
	for i := 0; i < argsNum-1; i++ {
		controllerParamsType = append(controllerParamsType, t.In(argsNum-1).Kind())
	}
	// 下方函数是运行时
	return func(c *gin.Context) {
		// 在框架初始化的时候通过反射获取类型同时注册路由，这样就不需要在controller里每次获取参数映射，而变成了函数参数
		var values []reflect.Value
		values = append(values, reflect.ValueOf(c))
		for index, param := range c.Params {
			var arg any
			if controllerParamsType[index] == reflect.Int {
				intId, err := strconv.Atoi(param.Value)
				if err != nil {
					MakeErrorResponse(c, errors2.Wrapper(errors2.ErrRouteParamInvalid, fmt.Sprintf(" controller path:%s url path:%s", param.Key, param.Value)))
					return
				}
				if controllerParamsType[index] == reflect.Int {
					arg = intId
				}
			} else {
				arg = param.Value
			}
			values = append(values, reflect.ValueOf(arg))
		}
		if argsNum > 1 && t.In(argsNum-1).Kind() == reflect.Pointer {
			bodyStruct = reflect.New(t.In(1).Elem()).Interface()
			param, err := a.parseBodyToJsonStruct(c, bodyStruct)
			if err != nil {
				return
			}
			values = append(values, reflect.ValueOf(param))
		}
		values = f.Call(values) // 执行控制器函数
		err := values[1].Interface()
		if err != nil {
			errInterface := err.(error)
			MakeErrorResponse(c, errInterface)
			return
		}
		data := values[0].Interface()
		MakeSuccessResponse(c, data)
	}
}

func (a *App) validateController(controller any) error {
	t := reflect.TypeOf(controller)
	if t.Kind() != reflect.Func {
		return nil
	}
	name := util.GetFunctionName(controller)
	if t.Kind() != reflect.Func {
		return errors.New(name + " controller must be func")
	}
	i := t.NumIn()
	if i < 1 {
		return errors.New(name + " controller must have at least one argument")
	}
	if t.In(0).Kind() != reflect.Ptr {
		return errors.New(name + " controller first argument must be gin context pointer")
	}
	if i > 1 && t.In(i-1).Kind() == reflect.Struct {
		return errors.New(name + " controller last argument must be param pointer")
	}
	// 校验 controller 出参定义
	if t.NumOut() != 2 {
		return errors.New(fmt.Sprintf("controller %s output args error", name))
	}
	if t.Out(0).Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("controller %s output first arg not ptr", name))
	}
	if t.Out(1).Kind() != reflect.Interface {
		return errors.New(fmt.Sprintf("controller %s output second arg not interface error", name))
	}
	if _, ok := reflect.New(t.Out(1)).Interface().(*error); !ok {
		return errors.New(fmt.Sprintf("controller %s output second arg not error", name))
	}
	return nil
}

func (a *App) parseUrlParams(c *gin.Context) []any {
	params := make([]any, 0)
	for _, param := range c.Params {
		params = append(params, param.Value)
	}
	return params
}

func (a *App) parseBodyToJsonStruct(c *gin.Context, reqStruct any) (any, error) {
	if err := c.ShouldBindJSON(&reqStruct); err != nil {
		errMsg := err.Error()
		if errMsg == "EOF" {
			errMsg = "empty json body"
		}
		traceId, _ := c.Request.Context().Value("traceId").(string)
		spanId, _ := c.Request.Context().Value("spanId").(string)
		c.JSON(400, gin.H{
			"traceId": traceId,
			"spanId":  spanId,
			"code":    400,
			"message": errMsg,
			"data":    nil,
		})
		return nil, err
	}
	return reqStruct, nil
}

func MakeSuccessResponse(c *gin.Context, data any) {
	traceId, _ := c.Request.Context().Value("traceId").(string)
	spanId, _ := c.Request.Context().Value("spanId").(string)
	response := entity.Response[any]{
		TraceId: traceId,
		SpanId:  spanId,
		Data:    data,
		Code:    0,
	}
	c.JSON(200, response)
}

func MakeErrorResponse(c *gin.Context, err error) {
	traceId, _ := c.Request.Context().Value("traceId").(string)
	spanId, _ := c.Request.Context().Value("spanId").(string)
	code := errors2.InternalServerErrorCode
	var message string
	httpStatus := http.StatusInternalServerError
	for _, e := range errors2.GetCollectErr() {
		if errors.Is(err, e) {
			code = errors2.GetErrCodeByErr(e)
			httpStatus = errors2.GetHttpStatusByErr(e)
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code = errors2.DataNotFound
		message = "数据不存在"
		httpStatus = http.StatusNotFound
	} else {
		message = err.Error()
	}
	c.JSON(httpStatus, gin.H{
		"traceId": traceId,
		"spanId":  spanId,
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// NewApp .
func NewApp(options *ServerOptions, middleware ...gin.HandlerFunc) *App {
	if app == nil {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000000",
		})
		logrus.SetOutput(os.Stdout)
		// disable gin log
		if options.Env == "production" {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}
		gin.DefaultWriter = io.Discard
		app = &App{
			Engine:        gin.Default(),
			command:       c,
			serverOptions: options,
		}
		app.Use(middleware...)
	}
	return app
}

// GetApp .
func GetApp() *App {
	return app
}

func (a *App) Running(ctx context.Context) error {
	log.Info(ctx, "run", map[string]any{
		"host": a.serverOptions.Host,
		"port": a.serverOptions.Port,
	})
	go func() {
		if err := a.command.Run(ctx); err != nil {
			log.Error(ctx, "command-error", err, nil)
		}
	}()
	if err := a.Run(fmt.Sprintf("%s:%d", a.serverOptions.Host, a.serverOptions.Port)); err != nil {
		panic(err)
	}
	return nil
}
