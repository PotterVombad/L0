package api

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/PotterVombad/L0/internal/db"
	"github.com/gin-gonic/gin"
)

type API struct {
	storage db.Store
}

const indexHTML = "index.html"

func (a API) search(c *gin.Context) {
	// TODO: just id
	// TODO: get
	uid := c.PostForm("uid")

	order, err := a.storage.Get(c.Request.Context(), uid)
	if err != nil {
		c.HTML(
			http.StatusInternalServerError,
			indexHTML,
			gin.H{"request": "Заказ не найден"},
		)
		return
	}

	bb, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		c.HTML(
			http.StatusInternalServerError,
			indexHTML,
			gin.H{"request": "Что то пошло не так"},
		)
		return
	}

	c.HTML(
		http.StatusOK, 
		indexHTML, 
		gin.H{"request": template.HTML("<pre>" + string(bb) + "</pre>")},
	)
}

func (a API) index(c *gin.Context) {
	c.HTML(http.StatusOK, indexHTML, gin.H{})
}

func (a API) Run() error {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	
	r.GET("/index", a.index)
	r.POST("/search", a.search)

	if err := r.Run(); err != nil {
		return err
	}
	return nil
}

func New(storage db.Store) API {
	return API{
		storage: storage,
	}
}
