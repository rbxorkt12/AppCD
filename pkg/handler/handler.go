package handler

import (
	"github.com/labstack/echo"
	"log"
	"net/http"

)

type Optionhandler struct{
	Change Changeoption
}

type Changeoption string

const YES Changeoption = "yes"
const No Changeoption = "no"
const Wait Changeoption = "wait"



func (o *Optionhandler) Sethandler(group *echo.Group) {
	group.POST("/yes", o.Yes)
	group.POST("/no", o.No)
}

func (o *Optionhandler) Yes(context echo.Context) error{
	log.Println("Set cluster use recent commit")
	//*o 로 접근하지않아도 되는건가?
	o.Change=YES
	return context.JSONPretty(http.StatusOK, "Change option is true", "    ")
}

func (o *Optionhandler) No(context echo.Context) error{
	log.Println("Don't Set cluster use recent commit")
	o.Change=No
	return context.JSONPretty(http.StatusOK, "Change option is false", "    ")
}

