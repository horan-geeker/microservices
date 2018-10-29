package Controller

import (
    "gopkg.in/go-playground/validator.v9"
    "microservices/Request"
    "log"
)

var validate *validator.Validate

func Test() {
    validate = validator.New()
    item := &Request.AddItem{
        Name: "",
    }
    err := validate.Struct(item)
    if err != nil {
        if _, ok := err.(*validator.InvalidValidationError); ok {
            log.Println(err)
            return
        }
        for _, err := range err.(validator.ValidationErrors) {
            log.Println(err.Namespace())
            log.Println(err.Field())
            log.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
            log.Println(err.StructField())     // by passing alt name to ReportError like below
            log.Println(err.Tag())
            log.Println(err.ActualTag())
            log.Println(err.Kind())
            log.Println(err.Type())
            log.Println(err.Value())
            log.Println(err.Param())
            log.Println(err)
        }
        return
    }
    log.Println("no error")
}