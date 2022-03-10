package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/machinebox/graphql"
)

// "clientCapabilities": [
// {
//   "key": "ads",
//   "value": [
//     "non_ssai"
//   ]
// },

// "xHsUsertoken":"token",
// "xCountryCode": "in",
// "xHsRequestId": "beb5a558-1335-4473-80af-26bb8d4c2529"

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
	}

	err := c.gqlClient.Run(ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func testgabs() {
	//jsonObj := gabs.Wrap("{}")
	jsonObj := gabs.New()
	// or gabs.Wrap(jsonObject) to work on an existing map[string]interface{}
	jsonObj.SetP("token", "hsRequest.xHsUsertoken")
	jsonObj.SetP("in", "hsRequest.xCountryCode")
	jsonObj.SetP("id", "hsRequest.xHsRequestId")
	tmp := gabs.New()
	tmp.Set("ads", "keys")
	tmp.Set([]string{"non_ssai", "value"}, "value")
	jsonObj.SetP([]interface{}{tmp}, "clientCapabilities")
	jsonObj.Set(10, "outter", "inner", "value")
	jsonObj.SetP(20, "outter.inner.value2")
	jsonObj.Set(30, "outter", "inner2", "value3")

	// fmt.Println(jsonObj.Data())
	fmt.Println(jsonObj.String())
	m := map[string]interface{}{
		"playbackRequestBody.header.xHsUsertoken": "token",
		"playbackRequestBody.header.xCountryCode": "in",
		"playbackRequestBody.header.xHsRequestId": "idname",
		"playbackRequestBody.contentId":           "1260023395",
	}
	d := writeToJson(m)
	bye, _ := json.Marshal(d)
	fmt.Println(string(bye))
	jsonObj = gabs.Wrap(map[string]interface{}{
		"header": map[string]interface{}{
			"key": "values",
		},
		"query_param": map[string]interface{}{
			"contentId": "1260023395",
			"hello":     1,
			"world":     []int32{2, 3},
		},
	})
	fmt.Println(jsonObj.Path("header.key"))
	fmt.Println(jsonObj.Path("query_param.contentId"))
}
func main() {
	testgabs()
	// test_variables()
}

func getparammap(pathVariablesMap map[string]string) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range pathVariablesMap {
		value := queryParam(v)
		m[k] = value
	}

	return m
}
func queryParam(key string) interface{} {
	//variables from request (header or query string)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1bV9hY2Nlc3MiLCJleHAiOjE2MzgyNTk0NzAsImlhdCI6MTYzODE3MzA3MCwiaXNzIjoiVFMiLCJqdGkiOiI0OWQ3NjE1NzM1YTM0YTc2OTljNjQwYzhhOGU3NjE3NCIsInN1YiI6IntcImhJZFwiOlwiNTZkNjhmYzZmZWNlNGYwMzg0OGE3NjRmZjIyYTM4OWRcIixcInBJZFwiOlwiYjI1NDZlOGJmNmIzNGEzYThmZjYzZDk3Y2ZiYTM2MDlcIixcImVtYWlsXCI6XCJodWFuLmxpdUBob3RzdGFyLmNvbVwiLFwiaXBcIjpcIjU5LjE0NC4xMDkuMTEwXCIsXCJjb3VudHJ5Q29kZVwiOlwiaW5cIixcImN1c3RvbWVyVHlwZVwiOlwibnVcIixcInR5cGVcIjpcImVtYWlsXCIsXCJpc0VtYWlsVmVyaWZpZWRcIjpmYWxzZSxcImlzUGhvbmVWZXJpZmllZFwiOmZhbHNlLFwiZGV2aWNlSWRcIjpcIjJkZTgzZGVlLTE0YzctNDA5Mi05NTBiLWUzNjI0MWJjMDNjYlwiLFwicHJvZmlsZVwiOlwiQURVTFRcIixcInZlcnNpb25cIjpcInYyXCIsXCJzdWJzY3JpcHRpb25zXCI6e1wiaW5cIjp7XCJIb3RzdGFyUHJlbWl1bVwiOntcInN0YXR1c1wiOlwiU1wiLFwiZXhwaXJ5XCI6XCIyMDIyLTA3LTMwVDE4OjMwOjAwLjAwMFpcIixcInNob3dBZHNcIjpcIjBcIixcImNudFwiOlwiMVwifX19LFwiZW50XCI6XCJDc0VCQ2dVS0F3b0JCUkszQVJJSFlXNWtjbTlwWkJJRGFXOXpFZ2xoYm1SeWIybGtkSFlTQm1acGNtVjBkaElIWVhCd2JHVjBkaElFY205cmRSSURkMlZpRWdSdGQyVmlFZ2QwYVhwbGJuUjJFZ1YzWldKdmN4SUdhbWx2YzNSaUVncGphSEp2YldWallYTjBFZ1IwZG05ekVnUndZM1IyRWdOcWFXOFNCMnBwYnkxc2VXWWFBbk5rR2dKb1pCb0RabWhrR2dJMGF5SURjMlJ5SWdWb1pISXhNQ0lMWkc5c1lubDJhWE5wYjI0cUJuTjBaWEpsYnlvSVpHOXNZbmsxTGpGWUFRb0xFZ2tJWkRnRVFBRlE4QkFLSWdvYUNnNFNCVFUxT0RNMkVnVTJOREEwT1FvSUlnWm1hWEpsZEhZU0JEaGtXQUVLd1FFS0JRb0RDZ0VBRXJjQkVnZGhibVJ5YjJsa0VnTnBiM01TQ1dGdVpISnZhV1IwZGhJR1ptbHlaWFIyRWdkaGNIQnNaWFIyRWdSeWIydDFFZ04zWldJU0JHMTNaV0lTQjNScGVtVnVkSFlTQlhkbFltOXpFZ1pxYVc5emRHSVNDbU5vY205dFpXTmhjM1FTQkhSMmIzTVNCSEJqZEhZU0EycHBieElIYW1sdkxXeDVaaG9DYzJRYUFtaGtHZ05tYUdRYUFqUnJJZ056WkhJaUJXaGtjakV3SWd0a2IyeGllWFpwYzJsdmJpb0djM1JsY21WdktnaGtiMnhpZVRVdU1WZ0JFZ2NRd0tEOWdxVXdcIixcImlzc3VlZEF0XCI6MTYzODE3MzA3MDUzN30iLCJ2ZXJzaW9uIjoiMV8wIn0.RGgKTgzGOxyKY0Vb7Uy-n4CxhJdQ6gE2UCH_GmCbHd4"
	country := "in"
	platform := "android"
	contentId := "1260023395"
	lpv := "eng"

	var paramsMap = map[string]interface{}{
		"token":     token,
		"country":   country,
		"platform":  platform,
		"contentId": contentId,
		"lpv":       lpv,
	}
	if v, ok := paramsMap[key]; ok {
		return v
	}
	return ""
}

func test_variables() {
	endpoint := "https://hsx-graphql-gateway.pp.hotstar-labs.com/"
	cli := &client{
		graphql.NewClient(endpoint),
	}
	ctx := context.TODO()
	params := map[string]interface{}{}
	// path from layout service, variables are defined in GraphQL Query
	pathVariablesMap := map[string]string{
		"hsRequest.token":        "token",
		"hsRequest.countryCode":  "country",
		"hsRequest.platformCode": "platform",
		"contentId":              "contentId",
		"lpv":                    "lpv",
	}
	m := getparammap(pathVariablesMap)
	d := writeToJson(m)
	params = d.Data().(map[string]interface{})
	fmt.Println(params)
	resp, err := cli.Query(ctx, query, params)
	if err != nil {
		log.Fatalf("err:%s", err)
	}
	for k, v := range resp {
		log.Printf("resp:%s, %s", k, v)
	}

}

// layout service query variables definition
// hsRequest.token:        token
// hsRequest.countryCode:  country
// hsRequest.platformCode: platform
// contentId: contentId
// lpv: lpv

var query = `query($hsRequest: HSRequest!, $contentId: String!,$lpv: String){
	content_detail:fetchContentDetail(
	hsRequest: $hsRequest 
	contentId: $contentId
	lpv: $lpv 
      ){
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
      } 
	}
`
