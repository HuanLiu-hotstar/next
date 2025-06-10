package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/diegosz/gqlformatter"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func formatQuery(s string) string {

	// s := "query{products(where:{and:{id:{gte:20} label:{eq:$label}}}){name    id  price }}"
	q, err := gqlformatter.FormatQuery(s)

	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(q)
	return q
}

type GQLQuery struct {
	WidgetName    string
	QueryPath     string
	Query         string
	FormatedQuery string
}

func getWidgetName(queryPath string) string {
	array := strings.Split(queryPath, "plugins")
	array = strings.Split(array[1], "/")
	return array[1] + "/" + array[2]

}

func getGQLQueryFilelist(root string) []GQLQuery {
	queryPathList := []GQLQuery{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}
		if strings.HasSuffix(path, ".gql") {
			fmt.Printf("Visited: %s\n", path)
			queryPathList = append(queryPathList, GQLQuery{
				WidgetName: getWidgetName(path),
				QueryPath:  path,
			})
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
	}
	return queryPathList
}

func readGQLQuery(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("err:%s, file:%s", err, path)
	}
	return string(content)
}

func main() {

	root := "/Users/huanliu/hotstar/hs-core-widget-binder-query/domains/consumption/plugins/PlayerWidgetInstance"
	gqlQueries := getGQLQueryFilelist(root)
	// formatedQueryList := []string{}
	for i, queryFile := range gqlQueries {
		query := readGQLQuery(queryFile.QueryPath)
		gqlQueries[i].Query = query
		fmt.Printf("GQLQuery:%v\n", queryFile)
		formatedQuery := formatQuery(query)
		if len(formatedQuery) != 0 {
			// formatedQueryList = append(formatedQueryList, formatedQuery)
			gqlQueries[i].FormatedQuery = formatedQuery
		}
	}
	for _, query := range gqlQueries {
		fmt.Printf("GQLQuery:%v\n", query)
	}
	main1()
}

func main1() {

	schema := `
		type Query {
			user(id: ID!): User
		}
		type User {
			id: ID!
			name: String!
			email: String
			posts: [Post!]!
		}
		type Post {
			title: String!
		}
	`

	// Load the schema
	schemaDoc, err := gqlparser.LoadSchema(&ast.Source{Input: schema})
	if err != nil {
		panic(err)
	}
	query := `
		query GetUser($id: ID!) {
			user(id: $id) {
				id
				name
				email
				posts {
					title
				}
			}
		}
	`

	// Parse the query
	// source := &ast.Source{Input: query}
	parsed, errList := gqlparser.LoadQuery(schemaDoc, query)
	if len(errList) != 0 {
		panic("error list")
	}

	// Print the operation type (query/mutation/subscription)
	fmt.Println("Operation Type:", parsed.Operations[0].Operation)

	// Print the operation name
	fmt.Println("Operation Name:", parsed.Operations[0].Name)

	// Print variables
	for _, variable := range parsed.Operations[0].VariableDefinitions {
		fmt.Printf("Variable: %s, Type: %s\n", variable.Variable, variable.Type.String())
	}

	// Print selected fields
	for _, selection := range parsed.Operations[0].SelectionSet {
		if field, ok := selection.(*ast.Field); ok {
			fmt.Println("Field:", field.Name)
			if field.SelectionSet != nil {
				for _, subSelection := range field.SelectionSet {
					if subField, ok := subSelection.(*ast.Field); ok {
						fmt.Println("  Subfield:", subField.Name)
					}
				}
			}
		}
	}
	fmt.Println(ast.Dump(parsed))
	fmt.Println(parsed.Operations[0].VariableDefinitions)
	ast.StringValue()
}
