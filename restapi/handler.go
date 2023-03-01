package restapi

import (
	"fmt"
	"immudblog/model"
	"immudblog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setAPIResponse(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, model.APIResponse{Code: code, Message: msg})
}

// Controller example
type Controller struct {
}

// NewController example
func NewController() *Controller {
	return &Controller{}
}

// GetLogs Getting logs
// @Summary      Get stored logs
// @Description  Get stored logs
// @Tags         logs
// @Accept       json
// @Produce      json
// @Security	 BasicAuth
// @Param        count    query    int  false  "Max number of returned logs. Ordered by id desc."
// @Param        application    query    string  false  "Filter for application"
// @Success      200  {array}   model.Log
// @Failure      400  {object}  model.APIResponse
// @Failure      500  {object}  model.APIResponse
// @Router       /logs [get]
func (c *Controller) GetLogs(ctx *gin.Context) {
	var cou uint64
	cou = 100
	couStr, found := ctx.GetQuery("count")
	if found {
		var err error
		cou, err = strconv.ParseUint(couStr, 10, 0)
		if err != nil {
			setAPIResponse(ctx, http.StatusBadRequest, "invalid 'count' query param")
			return
		}
	}
	app, found := ctx.GetQuery("application")
	if !found {
		app = ""
	}

	ret, err := service.GetLogs(cou, app)
	if err != nil {
		setAPIResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Internal error. Error: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, ret)
}

// AddLogs Add log(s) to the immudb
// @Summary      Add log(s) to the immudb
// @Description  Add log(s) to the immudb
// @Tags         logs
// @Accept       json
// @Produce      json
// @Security	 BasicAuth
// @Param 		request  body  []model.Log true "Logs to be storeds"
// @Success      200  {object}  model.APIResponse
// @Failure      400  {object}  model.APIResponse
// @Failure      500  {object}  model.APIResponse
// @Router       /logs [post]
func (c *Controller) AddLogs(ctx *gin.Context) {
	var logs []model.Log

	// bind the received JSON data to log
	if err := ctx.BindJSON(&logs); err != nil {
		setAPIResponse(ctx, http.StatusBadRequest, "invalid Log data")
		return
	}
	if err := service.AddLogs(logs); err != nil {
		setAPIResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Internal error. Error: %v", err))
		return
	}
	setAPIResponse(ctx, http.StatusOK, fmt.Sprintf("stored %d records", len(logs)))
}

// GetLogs Getting logs
// @Summary      Get number of stored logs
// @Description  Get number of stored logs
// @Tags         logs
// @Accept       json
// @Produce      json
// @Security	 BasicAuth
// @Success      200  {object}  model.GetLogsCountResponse
// @Failure      400  {object}  model.APIResponse
// @Failure      500  {object}  model.APIResponse
// @Router       /logs/count [get]
func (c *Controller) GetLogsCount(ctx *gin.Context) {

	ret, err := service.CountLogs()
	if err != nil {
		setAPIResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Internal error. Error: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, model.GetLogsCountResponse{Count: ret})
}
