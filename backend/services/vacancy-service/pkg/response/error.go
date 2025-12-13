package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	domainErrors "vacancy-service/internal/domain/errors"
)

func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func HandleError(c *gin.Context, err error, logger *zap.Logger) {
	logger.Error("handler error", zap.Error(err))

	switch err {
	case domainErrors.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case domainErrors.ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}

}
