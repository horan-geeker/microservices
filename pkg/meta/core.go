package meta

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	errors2 "microservices/errors"
	"microservices/internal/entity"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
)

var app *Engine

// Engine 包了一层 gin 内核
type Engine struct {
	*gin.Engine
}

type customHandler any

// POST is a shortcut for router.Handle("POST", path, handlers).
func (e *Engine) POST(relativePath string, handlers ...customHandler) gin.IRoutes {
	return e.Handle(http.MethodPost, relativePath, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handlers).
func (e *Engine) GET(relativePath string, handlers ...customHandler) gin.IRoutes {
	return e.Handle(http.MethodGet, relativePath, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handlers).
func (e *Engine) DELETE(relativePath string, handlers ...customHandler) gin.IRoutes {
	return e.Handle(http.MethodDelete, relativePath, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handlers).
func (e *Engine) PATCH(relativePath string, handlers ...customHandler) gin.IRoutes {
	return e.Handle(http.MethodPatch, relativePath, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handlers).
func (e *Engine) PUT(relativePath string, handlers ...customHandler) gin.IRoutes {
	return e.Handle(http.MethodPut, relativePath, handlers...)
}

func (e *Engine) Handle(httpMethod, relativePath string, customHandler ...customHandler) gin.IRoutes {
	handlers := make([]gin.HandlerFunc, 0)
	for _, handler := range customHandler {
		if err := e.validateController(handler); err != nil {
			panic(err)
		}
		handlers = append(handlers, e.wrapperGin(handler))
	}
	return e.Engine.Handle(httpMethod, relativePath, handlers...)
}

func (e *Engine) wrapperGin(handle any) gin.HandlerFunc {
	t := reflect.TypeOf(handle)
	f := reflect.ValueOf(handle)
	argsNum := t.NumIn()
	var urlParamsType []reflect.Kind
	var bodyStruct any
	// 控制器最后一个参数是结构体解析的 body
	for i := 0; i < argsNum-1; i++ {
		urlParamsType = append(urlParamsType, t.In(argsNum-1).Kind())
	}
	if argsNum > 1 && t.In(argsNum-1).Kind() == reflect.Pointer {
		bodyStruct = reflect.New(t.In(1).Elem()).Interface()
	}
	// 下方函数是运行时
	return func(c *gin.Context) {
		// 在框架初始化的时候通过反射获取类型同时注册路由，这样就不需要在controller里每次获取参数映射，而变成了函数参数
		var values []reflect.Value
		values = append(values, reflect.ValueOf(c))
		urlParams := parseUrlParams(c)
		if len(urlParams) > 0 {
			for index, urlParam := range urlParams {
				arg := urlParam
				var err error
				if urlParamsType[index] == reflect.Int {
					arg, err = strconv.Atoi(urlParam.(string))
					if err != nil {
						panic(err)
					}
				}
				values = append(values, reflect.ValueOf(arg))
			}
		}
		if bodyStruct != nil {
			param, err := parseBodyToJsonStruct(c, bodyStruct)
			if err != nil {
				return
			}
			values = append(values, reflect.ValueOf(param))
		}
		values = f.Call(values) // 执行控制器函数
		err := values[1].Interface()
		if err != nil {
			e := err.(error)
			makeErrorResponse(c, e)
			return
		}
		data := values[0].Interface().(map[string]any)
		makeSuccessResponse(c, data)
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (e *Engine) validateController(controller any) error {
	t := reflect.TypeOf(controller)
	if t.Kind() != reflect.Func {
		return nil
	}
	name := getFunctionName(controller)
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
	if t.Out(0).Kind() != reflect.Map {
		return errors.New(fmt.Sprintf("controller %s output first arg not map", name))
	}
	if t.Out(1).Kind() != reflect.Interface {
		return errors.New(fmt.Sprintf("controller %s output second arg not interface error", name))
	}
	if _, ok := reflect.New(t.Out(1)).Interface().(*error); !ok {
		return errors.New(fmt.Sprintf("controller %s output second arg not error", name))
	}
	return nil
}

func parseUrlParams(c *gin.Context) []any {
	params := make([]any, 0)
	for _, param := range c.Params {
		params = append(params, param.Value)
	}
	return params
}

func parseBodyToJsonStruct(c *gin.Context, reqStruct any) (any, error) {
	if err := c.ShouldBindJSON(&reqStruct); err != nil {
		requestId, _ := c.Request.Context().Value("requestId").(string)
		c.JSON(400, gin.H{
			"request_id": requestId,
			"code":       400,
			"message":    err.Error(),
			"data":       nil,
		})
		return nil, err
	}
	return reqStruct, nil
}

func makeSuccessResponse(c *gin.Context, data map[string]any) {
	requestId, _ := c.Request.Context().Value("requestId").(string)
	response := entity.Response{
		RequestId: requestId,
		Data:      data,
		Code:      0,
	}
	c.JSON(200, response)
}

func makeErrorResponse(c *gin.Context, err error) {
	requestId, _ := c.Request.Context().Value("requestId").(string)
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
		"request_id": requestId,
		"code":       code,
		"message":    message,
		"data":       nil,
	})
}

// GetEnginInstance .
func GetEnginInstance() *Engine {
	if app == nil {
		app = &Engine{
			Engine: gin.Default(),
		}
		return app
	}
	return app
}