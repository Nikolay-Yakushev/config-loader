package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Adapter) getVariable(ctx *gin.Context) {
	var stenvType string
	param, ok := ctx.GetQuery("source")
	if ok != true {
		err := errors.New("no source query parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": ctx.Error(err)},
		)
		return
	}
	l := a.log.Sugar()
	newL := l.With("source", param)
	switch param {
		
		case "file":
			newL.Debugf("Query param `source` reading")
			stenvType = "file"
		case "env":
			newL.Debugf("Query param `source` reading")
			stenvType = "env"
		default:
			newL.Debugf("Neither `env` or `file` query parameter was provided")
			ctx.JSON(http.StatusBadRequest, gin.H{"invalid storage type": param})
			return
	}
	value, err := a.storage[stenvType].GetValue(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ctx.Error(err)})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"value": value},
	)
	return 
	 
}
