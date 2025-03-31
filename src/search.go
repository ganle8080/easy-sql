package easysql

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

// SearchSchema 通用查询schema结构体
type SearchSchema struct {
	TableName string      `json:"table_name"`
	Columns   []Column    `json:"columns"`
	Joins     []Join      `json:"joins"`
	Wheres    []Condition `json:"wheres"`
	Orders    []string    `json:"orders"`
	Groups    []string    `json:"groups"`
}

// Column 字段
type Column struct {
	Field   string `json:"field"`
	Alias   string `json:"alias"`
	Handler string `json:"handler"`
}

// Condition 条件
type Where struct {
	Field   string `json:"field"`
	Handler string `json:"handler"`
}

// Join 链接参数
type Join struct {
	LinkField      string `json:"link_field"`
	JoinTableName  string `json:"join_table_name"`
	JoinTableField string `json:"join_table_field"`
	JoinType       string `json:"join_type"`
}

type FieldHandler struct {
	Name        string
	HandlerName string
	MethodName  string
	Args        []string
}

type SearchData struct {
	SchemaName string      `json:"schema_name"`
	Fields     []string    `json:"fields"`
	Conditions []Condition `json:"conditions"`
	Orders     []string    `json:"orders"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
}

type Condition struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func Search(db *gorm.DB, searchData *SearchData, schema *SearchSchema) {

}

// loadSearchSchema 查询指定路径下的json
func loadSearchSchema(schemaName string) (SearchSchema, error) {
	// 拼接路径
	filePath := schemaName + "_search_schema.json"

	// 加载文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return SearchSchema{}, err
	}

	// 类型转换
	var searchSchema SearchSchema
	jsonErr := json.Unmarshal(data, &schemaName)
	if jsonErr != nil {
		return SearchSchema{}, jsonErr
	}

	return searchSchema, nil
}

func buildColumns(searchFields []string, columns []Column) (resultFields []string, handlers []FieldHandler, err error) {
	searchFieldStr := strings.Join(searchFields, ";")

	// 判断fieldStr是否为空
	isNull := len(searchFieldStr) == 0

	for _, c := range columns {
		if isNull || strings.Contains(searchFieldStr, c.Alias) {
			resultFields = append(resultFields, c.Field+" AS "+c.Alias)
		}

		if (isNull || strings.Contains(searchFieldStr, c.Alias)) && len(c.Handler) > 0 {
			handlerStructStrs := strings.Split(c.Handler, ";")
			handlerAndMethodStrs := strings.Split(handlerStructStrs[0], ".")
			// handlerStructStrs的长度必须大于等于2,handlerAndMethodStrs必须等于2
			if len(handlerStructStrs) >= 2 && len(handlerAndMethodStrs) == 2 {
				handlerName := handlerAndMethodStrs[0]
				handlerMethodName := handlerAndMethodStrs[1]
				handlerMethodArgs := handlerStructStrs[1:]
				handler := FieldHandler{
					Name:       handlerName,
					MethodName: handlerMethodName,
					Args:       handlerMethodArgs,
				}
				handlers = append(handlers, handler)
			} else {
				err = fmt.Errorf("schema error handler not conform to the standards ,column is %s .", c.Alias)
				return
			}
		}
	}
	return
}

func buildWheres(conditions []Condition, wheres []Where) (whereStrs []string, handlers []FieldHandler, err error) {
	wheresMap := map[string]Where{}
	for _, w := range wheres {
		wheresMap[w.Field] = w
	}

	for _, c := range conditions {

		where, ok := wheresMap[c.Name]

		if !ok {
			err = fmt.Errorf("%s not in schema ", c.Name)
			return
		}

		fmt.Printf("where: %v\n", where)

		switch c.Type {
		case "eq":
			str := c.Name + " = " + c.Value
			whereStrs = append(whereStrs, str)
		case "ne":
		case "lt":
		case "gt":
		case "like":
		case "null":
		}

	}
}
