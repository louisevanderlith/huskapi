package huskapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/louisevanderlith/husk"
)

func NewAPI(ctx husk.Ctxer, defaultPagesize int, middleware ...gin.HandlerFunc) *gin.Engine {
	result := gin.Default()

	tbls := husk.TableLayouts(ctx)
	for name, tbl := range tbls {
		lname := strings.ToLower(name)
		result.GET(fmt.Sprintf("/%s/:key", lname), viewAction(tbl))

		authed := result.Group(fmt.Sprintf("/%s", lname))
		authed.Use(middleware...)
		authed.POST("", createAction(tbl))
		authed.PUT("/:key", updateAction(tbl))
		authed.DELETE("/:key", deleteAction(tbl))

		result.GET(fmt.Sprintf("/%s", lname), getAction(tbl, defaultPagesize))
		result.GET(fmt.Sprintf("/search/%s/:pagesize/*hash", lname), searchAction(tbl, defaultPagesize))
	}

	return result
}

//empty search
func getAction(t husk.Tabler, pageSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		results := t.Find(1, pageSize, husk.Everything())
		c.JSON(http.StatusOK, results)
	}
}

func searchAction(t husk.Tabler, pageSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, size := GetPageData(c.Param("pagesize"))
		decoded, err := base64.StdEncoding.DecodeString(c.Param("hash"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		parmObj := reflect.New(t.Type()).Interface().(husk.Dataer)
		err = json.Unmarshal(decoded, parmObj)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		results := t.Find(page, size, husk.ByFields(parmObj))

		c.JSON(http.StatusOK, results)
	}
}

func viewAction(t husk.Tabler) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.Param("key")
		key, err := husk.ParseKey(k)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		rec, err := t.FindByKey(key)

		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.JSON(http.StatusOK, rec)
	}
}

func createAction(t husk.Tabler) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := reflect.New(t.Type()).Interface().(husk.Dataer)
		err := c.Bind(&body)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		cset := t.Create(body)

		if cset.Error != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = t.Save()

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, cset.Record)
	}
}

func updateAction(t husk.Tabler) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := husk.ParseKey(c.Param("key"))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		body := reflect.New(t.Type()).Interface().(husk.Dataer)
		err = c.Bind(body)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		obj, err := t.FindByKey(key)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = obj.Set(body)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = t.Update(obj)

		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		err = t.Save()

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func deleteAction(t husk.Tabler) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := c.Param("key")
		key, err := husk.ParseKey(k)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = t.Delete(key)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = t.Save()

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, "Completed")
	}
}

func GetPageData(pageData string) (int, int) {
	defaultPage := 1
	defaultSize := 10

	if len(pageData) < 2 {
		return defaultPage, defaultSize
	}

	pChar := []rune(pageData[:1])

	if len(pChar) != 1 {
		return defaultPage, defaultSize
	}

	page := int(pChar[0]) % 32
	pageSize, err := strconv.Atoi(pageData[1:])

	if err != nil {
		return defaultPage, defaultSize
	}

	return page, pageSize
}
