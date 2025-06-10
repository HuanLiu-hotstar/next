package main

import "github.com/gin-gonic/gin"

func main() {
    api := gin.Default()
    api.GET("/pages/:page_id", func(ctx *gin.Context) {

        pageId := ctx.Param("page_id")
        offset := ctx.Query("offset")
        size := ctx.Query("size")
        ctx.JSON(200, map[string]interface{}{
            "page_id": pageId,
            "offset":  offset,
            "size":    size,
        })
    })
    api.Run(":8080")
}
