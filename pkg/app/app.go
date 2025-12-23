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
	"microservices/pkg/ecode"
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
		if strings.Contains(name, "/controller.") {
			if err := a.validateController(handler); err != nil {
				panic(err)
			}
			ginHandlerFunc = a.wrapperGin(handler)
		} else {
			var ok bool
			ginHandlerFunc, ok = handler.(gin.HandlerFunc)
			if !ok {
				log.Error(context.Background(), "system-error", nil, map[string]any{
					"function": name,
				})
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

	// 分析控制器参数结构
	hasBodyParam := argsNum > 1 && t.In(argsNum-1).Kind() == reflect.Pointer
	pathParamStart := 1 // 第一个参数是 gin.Context，从第二个开始是路径参数
	pathParamEnd := argsNum
	if hasBodyParam {
		pathParamEnd = argsNum - 1 // 最后一个是请求体参数，排除
	}

	// 收集路径参数的类型信息
	var pathParamTypes []reflect.Kind
	for i := pathParamStart; i < pathParamEnd; i++ {
		pathParamTypes = append(pathParamTypes, t.In(i).Kind())
	}

	// 下方函数是运行时
	return func(c *gin.Context) {
		var values []reflect.Value
		values = append(values, reflect.ValueOf(c))

		// 处理路径参数：动态匹配路径参数数量
		pathParams := c.Params
		expectedPathParamCount := len(pathParamTypes)
		actualPathParamCount := len(pathParams)

		// 验证路径参数数量是否匹配
		if actualPathParamCount != expectedPathParamCount {
			MakeErrorResponse(c, ecode.ErrRouteParamInvalid.WithMessage(
				fmt.Sprintf("路径参数数量不匹配: 期望 %d 个，实际 %d 个", expectedPathParamCount, actualPathParamCount)))
			return
		}

		// 按顺序处理每个路径参数
		for i, param := range pathParams {
			var arg any
			paramType := pathParamTypes[i]

			if paramType == reflect.Int {
				intId, err := strconv.Atoi(param.Value)
				if err != nil {
					MakeErrorResponse(c, ecode.ErrRouteParamInvalid.WithMessage(
						fmt.Sprintf("路径参数 %s 无法转换为整数: %s", param.Key, param.Value)))
					return
				}
				arg = intId
			} else if paramType == reflect.String {
				arg = param.Value
			} else {
				// 默认按字符串处理
				arg = param.Value
			}
			values = append(values, reflect.ValueOf(arg))
		}

		// 处理请求体参数
		if hasBodyParam {
			bodyStruct := reflect.New(t.In(argsNum - 1).Elem()).Interface()
			param, err := a.parseBodyToJsonStruct(c, bodyStruct)
			if err != nil {
				return
			}
			values = append(values, reflect.ValueOf(param))
		}

		// 执行控制器函数
		results := f.Call(values)
		err := results[1].Interface()
		if err != nil {
			errInterface := err.(error)
			MakeErrorResponse(c, errInterface)
			return
		}
		data := results[0].Interface()
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

	argCount := t.NumIn()
	if argCount < 1 {
		return errors.New(name + " controller must have at least one argument (gin.Context)")
	}

	// 第一个参数必须是 *gin.Context
	if t.In(0) != reflect.TypeOf((*gin.Context)(nil)) {
		return errors.New(name + " controller first argument must be *gin.Context")
	}

	// 支持以下几种组合的路由
	//  | /courses/:id                              | func(c *gin.Context, id string)                             | /courses/123                        |
	//  | /courses/:id/students                     | func(c *gin.Context, id string, param *RequestStruct)       | /courses/123/students               |
	//  | /courses/:courseId/students/:studentId    | func(c *gin.Context, courseId string, studentId string)     | /courses/123/students/456           |
	//  | /api/:version/courses/:id/actions/:action | func(c *gin.Context, version string, id int, action string) | /api/v1/courses/123/actions/publish |
	// 验证路径参数类型（第2个到倒数第2个或最后一个参数）
	hasBodyParam := argCount > 1 && t.In(argCount-1).Kind() == reflect.Pointer
	pathParamEnd := argCount
	if hasBodyParam {
		pathParamEnd = argCount - 1
	}

	// 验证路径参数类型，支持 string 和 int 类型
	for i := 1; i < pathParamEnd; i++ {
		paramType := t.In(i).Kind()
		if paramType != reflect.String && paramType != reflect.Int {
			return errors.New(fmt.Sprintf("%s controller path parameter %d must be string or int, got %s",
				name, i, paramType.String()))
		}
	}

	// 如果有请求体参数，验证最后一个参数必须是结构体指针
	if hasBodyParam {
		if t.In(argCount-1).Kind() != reflect.Pointer || t.In(argCount-1).Elem().Kind() != reflect.Struct {
			return errors.New(name + " controller last argument must be struct pointer for request body")
		}
	}

	// 校验 controller 出参定义
	if t.NumOut() != 2 {
		return errors.New(fmt.Sprintf("controller %s must return exactly 2 values", name))
	}

	// 第一个返回值必须是指针类型（响应结构体指针）
	if t.Out(0).Kind() != reflect.Pointer {
		return errors.New(fmt.Sprintf("controller %s first return value must be pointer", name))
	}

	// 第二个返回值必须是 error 接口
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !t.Out(1).Implements(errorInterface) {
		return errors.New(fmt.Sprintf("controller %s second return value must implement error interface", name))
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
	var err error
	// 根据请求方法决定绑定方式
	if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
		// GET/DELETE 请求绑定查询参数
		err = c.ShouldBind(reqStruct)
	} else {
		// POST/PUT/PATCH 请求绑定 JSON 参数
		err = c.ShouldBind(reqStruct)
	}

	if err != nil {
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
	code := ecode.InternalServerErrorCode
	var message string
	httpStatus := http.StatusInternalServerError
	var e *ecode.CustomError
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
		httpStatus = e.HttpStatus
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		code = ecode.DataNotFound
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
