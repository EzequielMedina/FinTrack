package router

import (
    "github.com/fintrack/chatbot-service/internal/infrastructure/entrypoints/handlers"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(handler *handlers.ChatHandler) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    r.GET("/health", handler.Health)

    api := r.Group("/api/chat")
    {
        api.POST("/query", handler.Query)
        api.POST("/report/pdf", handler.ReportPDF)
        api.POST("/report/chart", handler.ReportChart)
    }
    return r
}