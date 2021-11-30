package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
)

const (
	defaultKey = "_jsonkey_"
	jsonKey    = "jsonkey"
	headersKey = "header"
	queryKey   = "query_param"
)

func writeToJson(m map[string]interface{}) *gabs.Container {
	jsonObj := gabs.New()

	for k, v := range m {
		jsonObj.SetP(v, k)

	}
	return jsonObj
}

type client struct {
	gqlClient *graphql.Client
}
type Api interface {
	Query(ctx context.Context, query string, params map[string]interface{}) (map[string]json.RawMessage, error)
}

func Init(endpoint string) Api {
	gqlClient := graphql.NewClient(endpoint)

	return &client{
		gqlClient: gqlClient,
	}
}
func (c client) Query(ctx context.Context, query string, params map[string]interface{}) (map[string]json.RawMessage, error) {

	var resp map[string]json.RawMessage

	req := graphql.NewRequest(query)
	for k, v := range params {
		req.Var(k, v)
		log.Printf("variables:%s,%+v", k, v)
	}

	err := c.gqlClient.Run(ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {

	endpoint := "https://hsx-graphql-gateway.pp.hotstar-labs.com/"
	cli := &client{
		graphql.NewClient(endpoint),
	}
	router := gin.Default()
	router.Use(emptyGabMiddle)
	router.Use(gabHeaderMiddle)
	router.Use(func(c *gin.Context) {
		jsonMap := map[string]interface{}{}
		queryMap := getQueryStringMap(c)
		log.Println("query_param_map:", queryMap)
		if v, ok := queryMap[defaultKey]; ok {
			err := json.Unmarshal([]byte(v), &jsonMap)
			if err != nil {
				c.AbortWithError(500, err)
			}
			log.Println("jsonMap:%+v", jsonMap)
		}
		constructMap(c, queryMap, jsonMap)
		c.Next()
	})
	router.GET("/hello/:page", func(c *gin.Context) {
		page := c.Params.ByName("page")
		query, paramMap := getLayout(c, page)
		log.Printf("query:%s,paramMap:%+v", query, paramMap)
		variables := getParamMap(c, paramMap)
		resp, err := cli.Query(c, query, variables)
		if err != nil {
			c.AbortWithError(500, err)
		}
		c.JSON(200, resp)

	})
	router.Run(":8080")
}
func emptyGabMiddle(c *gin.Context) {
	if _, ok := c.Value(jsonKey).(*gabs.Container); !ok {
		c.Set(jsonKey, gabs.New())
	}
	c.Next()
}

//example for passport
func passportGabMiddle(c *gin.Context) {

	var ga *gabs.Container
	if v, ok := c.Value(jsonKey).(*gabs.Container); !ok {
		ga = gabs.New()
	} else {
		ga = v
	}
	passport := c.Request.Header.Get("X-Hs-Passport")
	m := map[string]interface{}{}
	// we should decode passport,this may be done by Passport-SDK and unmashal it to map[string]interface{}
	// then ,place it in ga
	json.Unmarshal([]byte(passport), m)
	ga.Set(m, "passport")
	c.Next()
}
func gabHeaderMiddle(c *gin.Context) {
	var ga *gabs.Container
	if v, ok := c.Value(jsonKey).(*gabs.Container); !ok {
		ga = gabs.New()
	} else {
		ga = v
	}
	headerMap := map[string]interface{}{}
	for k := range c.Request.Header {
		headerMap[k] = c.Request.Header.Get(k)
	}
	ga.Set(headerMap, headersKey)
	c.Next()
}
func getMapInterface(m map[string]string) map[string]interface{} {
	r := make(map[string]interface{}, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r

}
func constructMap(ctx context.Context, queryMap map[string]string, jsonkeyMap map[string]interface{}) *gabs.Container {

	if _, ok := queryMap[defaultKey]; ok {
		delete(queryMap, defaultKey)
	}
	var ga *gabs.Container
	if v, ok := ctx.Value(jsonKey).(*gabs.Container); ok {
		ga = v
	} else {
		ga = gabs.New()
	}

	ga.Set(getMapInterface(queryMap), queryKey)
	ga.Set(jsonkeyMap, jsonKey)

	log.Printf("ga:%s", ga.String())
	return ga
}
func getHeader(ctx *gin.Context) map[string]string {
	temp := map[string]string{}
	for k, _ := range ctx.Request.Header {
		temp[k] = ctx.Request.Header.Get(k)
	}
	return temp
}
func getQueryStringMap(ctx *gin.Context) map[string]string {
	temp := map[string]string{}
	for k, v := range ctx.Request.URL.Query() {
		temp[k] = strings.Join(v, ",")
	}
	return temp
}
func getParamMap(ctx context.Context, m map[string]string) map[string]interface{} {
	jsonObj, ok := ctx.Value(jsonKey).(*gabs.Container)
	if !ok {
		return nil
	}
	log.Println("jsonObj", jsonObj.String())
	resultmap := map[string]interface{}{}
	for resultKey, pathKey := range m {
		value := jsonObj.Path(pathKey)
		log.Printf("resultKey:%s,pathkey:%s, value:%v\n", resultKey, pathKey, value)
		resultmap[resultKey] = value
	}
	variables := writeToJson(resultmap)

	return variables.Data().(map[string]interface{})
}

func getLayout(c *gin.Context, page string) (string, map[string]string) {
	pathVariablesMap := map[string]string{
		"hsRequest.token":        "jsonkey.token",
		"hsRequest.countryCode":  "header.Country",
		"hsRequest.platformCode": "query_param.platform",
		"contentId":              "query_param.content_id",
		"lpv":                    "header.Lpv",
	}
	return query, pathVariablesMap
}

var query = `
query($hsRequest: HSRequest!, $contentId: String!, $lpv: String) {
	fetchContentDetail(hsRequest: $hsRequest, contentId: $contentId, lpv: $lpv) {
	  contentMeta {
	    coreAttributes {
	      title
	      description
	      images {
		verticalImage
		horizontalImage
		titleImage
	      }
	      isPaid
	      languages {
		code
		displayName
	      }
	      genres
	      ageRating
	      studioName
	      socialEnabled
	    }
	  }
	  titleCutOut {
	    imageUrl
	  }
      
	  watchProgress {
	    contentId
	    resumeAt
	    duration
	    watchState
	    watchRatio
	    lastUpdated
	    watchTag
	  }
	  # watchListStatus
	}
      }
`

// query example:
// curl --location -v 'localhost:8080/hello/world?platform=ios&content_id=1260023395&_jsonkey_=\{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1bV9hY2Nlc3MiLCJleHAiOjE2MzgyNTk0NzAsImlhdCI6MTYzODE3MzA3MCwiaXNzIjoiVFMiLCJqdGkiOiI0OWQ3NjE1NzM1YTM0YTc2OTljNjQwYzhhOGU3NjE3NCIsInN1YiI6IntcImhJZFwiOlwiNTZkNjhmYzZmZWNlNGYwMzg0OGE3NjRmZjIyYTM4OWRcIixcInBJZFwiOlwiYjI1NDZlOGJmNmIzNGEzYThmZjYzZDk3Y2ZiYTM2MDlcIixcImVtYWlsXCI6XCJodWFuLmxpdUBob3RzdGFyLmNvbVwiLFwiaXBcIjpcIjU5LjE0NC4xMDkuMTEwXCIsXCJjb3VudHJ5Q29kZVwiOlwiaW5cIixcImN1c3RvbWVyVHlwZVwiOlwibnVcIixcInR5cGVcIjpcImVtYWlsXCIsXCJpc0VtYWlsVmVyaWZpZWRcIjpmYWxzZSxcImlzUGhvbmVWZXJpZmllZFwiOmZhbHNlLFwiZGV2aWNlSWRcIjpcIjJkZTgzZGVlLTE0YzctNDA5Mi05NTBiLWUzNjI0MWJjMDNjYlwiLFwicHJvZmlsZVwiOlwiQURVTFRcIixcInZlcnNpb25cIjpcInYyXCIsXCJzdWJzY3JpcHRpb25zXCI6e1wiaW5cIjp7XCJIb3RzdGFyUHJlbWl1bVwiOntcInN0YXR1c1wiOlwiU1wiLFwiZXhwaXJ5XCI6XCIyMDIyLTA3LTMwVDE4OjMwOjAwLjAwMFpcIixcInNob3dBZHNcIjpcIjBcIixcImNudFwiOlwiMVwifX19LFwiZW50XCI6XCJDc0VCQ2dVS0F3b0JCUkszQVJJSFlXNWtjbTlwWkJJRGFXOXpFZ2xoYm1SeWIybGtkSFlTQm1acGNtVjBkaElIWVhCd2JHVjBkaElFY205cmRSSURkMlZpRWdSdGQyVmlFZ2QwYVhwbGJuUjJFZ1YzWldKdmN4SUdhbWx2YzNSaUVncGphSEp2YldWallYTjBFZ1IwZG05ekVnUndZM1IyRWdOcWFXOFNCMnBwYnkxc2VXWWFBbk5rR2dKb1pCb0RabWhrR2dJMGF5SURjMlJ5SWdWb1pISXhNQ0lMWkc5c1lubDJhWE5wYjI0cUJuTjBaWEpsYnlvSVpHOXNZbmsxTGpGWUFRb0xFZ2tJWkRnRVFBRlE4QkFLSWdvYUNnNFNCVFUxT0RNMkVnVTJOREEwT1FvSUlnWm1hWEpsZEhZU0JEaGtXQUVLd1FFS0JRb0RDZ0VBRXJjQkVnZGhibVJ5YjJsa0VnTnBiM01TQ1dGdVpISnZhV1IwZGhJR1ptbHlaWFIyRWdkaGNIQnNaWFIyRWdSeWIydDFFZ04zWldJU0JHMTNaV0lTQjNScGVtVnVkSFlTQlhkbFltOXpFZ1pxYVc5emRHSVNDbU5vY205dFpXTmhjM1FTQkhSMmIzTVNCSEJqZEhZU0EycHBieElIYW1sdkxXeDVaaG9DYzJRYUFtaGtHZ05tYUdRYUFqUnJJZ056WkhJaUJXaGtjakV3SWd0a2IyeGllWFpwYzJsdmJpb0djM1JsY21WdktnaGtiMnhpZVRVdU1WZ0JFZ2NRd0tEOWdxVXdcIixcImlzc3VlZEF0XCI6MTYzODE3MzA3MDUzN30iLCJ2ZXJzaW9uIjoiMV8wIn0.RGgKTgzGOxyKY0Vb7Uy-n4CxhJdQ6gE2UCH_GmCbHd4"\}' -H 'country: in' -H 'lpv: eng'
