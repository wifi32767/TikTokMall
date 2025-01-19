package utils

import (
	"net/http"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
)

func ErrorHandle(c *gin.Context, err error) {
	if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
		c.JSON(int(bizErr.BizStatusCode()), gin.H{"error": bizErr.BizMessage()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
