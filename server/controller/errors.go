package controller

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/domain"
)

func writeDomainError(c *gin.Context, err error) {
	if err == nil {
		c.Status(nethttp.StatusInternalServerError)
		return
	}

	if de, ok := err.(*domain.Error); ok {
		status := mapKindToStatus(de.Kind)
		c.JSON(status, gin.H{
			"error":  de.Msg,
			"kind":   de.Kind,
			"field":  de.Field,
			"entity": de.Entity,
		})
		return
	}

	c.JSON(nethttp.StatusInternalServerError, gin.H{"error": err.Error()})
}

func mapKindToStatus(k domain.Kind) int {
	switch k {
	case domain.Required, domain.Invalid:
		return nethttp.StatusBadRequest
	case domain.NotFound:
		return nethttp.StatusNotFound
	case domain.AlreadyExists, domain.Conflict:
		return nethttp.StatusConflict
	case domain.Unauthorized:
		return nethttp.StatusUnauthorized
	case domain.Forbidden:
		return nethttp.StatusForbidden
	default:
		return nethttp.StatusInternalServerError
	}
}

func domainErrorFromValue(err domain.Error) *domain.Error {
	if err.Kind == "" && err.Msg == "" && err.Err == nil {
		return nil
	}
	return &err
}
