package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koron/go-dproxy"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() { http.Handle("/", new()) }

func new() http.Handler {
	r := gin.New()
	r.Use(logErr())
	r.POST("/json", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var v interface{}
		if err = json.Unmarshal(body, &v); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		p := dproxy.New(v)
		url, err := p.M("share").String()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		ctx := appengine.NewContext(c.Request)
		log.Infof(ctx, "speedtest: %s", body)
		c.String(http.StatusOK, url)
	})
	return r
}

func logErr() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ctx := appengine.NewContext(c.Request)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Errorf(ctx, "%s", e.Error())
			}
		}
	}
}
