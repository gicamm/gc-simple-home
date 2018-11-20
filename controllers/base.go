package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"runtime/debug"
)

//BaseController is used as interceptor

type BaseController struct {
	beego.Controller
}

/*func (this *BaseController) Prepare() {
	utility.PccLogger.Info(this.Ctx.Request.RequestURI + " Has been called")
}*/

/*func (this *BaseController) Finish() {
	utility.PccLogger.Info(this.Ctx.Request.RequestURI + " Has been finished")
}*/

func (this *BaseController) Options() {
	this.ServeJSON()
}

//Generic Panic handler for all requests.
func (this *BaseController) handlePanic() {
	if r := recover(); r != nil {
		debug.PrintStack()
		status := this.Ctx.ResponseWriter.Status
		if status == 0 {
			status = 500
			this.Ctx.ResponseWriter.WriteHeader(status)
		}
		errorMessage := fmt.Sprintf("%v", r)
		this.Data["json"] = CreateResponse(status, errorMessage, "InternalServerError", nil)
		this.ServeJSON()
	}
}
