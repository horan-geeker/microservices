package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
	"microservices/internal/pkg/options"
	"microservices/internal/store/mysql"
	"microservices/internal/store/redis"
	"reflect"
	"runtime"
	"strconv"
)

var routes []router
var dataFactory = mysql.GetMysqlInstance(options.NewMySQLOptions())
var cacheFactory = redis.GetRedisInstance(options.NewRedisOptions())

type router struct {
	Method        string
	Path          string
	BeforeHandler []any
	Func          any
	AfterHandler  []any
}

// Register router register
func Register(r *gin.Engine) {
	for _, route := range routes {
		if err := validateController(route.Func); err != nil {
			panic(err)
		}
		r.Handle(route.Method, route.Path, wrapperGin(route))
	}
}

func wrapperGin(route router) gin.HandlerFunc {
	t := reflect.TypeOf(route.Func)
	f := reflect.ValueOf(route.Func)
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
		code := values[1].Int()
		err := values[2].Interface()
		if err != nil {
			e := err.(error)
			makeErrorResponse(c, code, e.Error())
			return
		}
		data := values[0].Interface().(map[string]any)
		makeSuccessResponse(c, data, code)
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func validateController(controller any) error {
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
	o := t.NumOut()
	if o != 3 {
		return errors.New(name + " controller must return 3 values")
	}
	if t.Out(0).Kind() != reflect.Map {
		return errors.New(name + " controller first return value must be map[string]any")
	}
	if t.Out(1).Kind() != reflect.Int {
		return errors.New(name + " controller second return value must be int")
	}
	if t.Out(2).Kind() != reflect.Interface {
		return errors.New(name + " controller third return value must be error")
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

func makeSuccessResponse(c *gin.Context, data map[string]any, code int64) {
	requestId, _ := c.Request.Context().Value("requestId").(string)
	response := entity.Response{
		RequestId: requestId,
		Data:      data,
		Code:      code,
	}
	c.JSON(200, response)
}

func makeErrorResponse(c *gin.Context, code int64, message string) {
	requestId, _ := c.Request.Context().Value("requestId").(string)
	c.JSON(500, gin.H{
		"request_id": requestId,
		"code":       code,
		"message":    message,
		"data":       nil,
	})
}
