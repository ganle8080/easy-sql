package easysql

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ganle8080/easysql/config/handler"
	"gorm.io/gorm"
)

// SearchSchema 通用查询schema结构体
type SearchSchema struct {
	TableName string      `json:"table_name"`
	Columns   []Column    `json:"columns"`
	Joins     []Join      `json:"joins"`
	Wheres    []Condition `json:"wheres"`
	Orders    []string    `json:"orders"`
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
	Field       string
	HandlerName string
	MethodName  string
	Args        []interface{}
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
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
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

func buildColumns(searchFields []string, columns []Column) (columnStr string, handlers []FieldHandler, err error) {

	var resultFields []string

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
				handlerMethodArgs := []interface{}{}
				if len(handlerStructStrs[1:]) > 0 {
					for _, v := range handlerStructStrs[1:] {
						handlerMethodArgs = append(handlerMethodArgs, v)
					}
				}
				handler := FieldHandler{
					Field:       c.Alias,
					HandlerName: handlerName,
					MethodName:  handlerMethodName,
					Args:        handlerMethodArgs,
				}
				handlers = append(handlers, handler)
			} else {
				err = fmt.Errorf("schema error handler not conform to the standards ,column is %s .", c.Alias)
				return
			}
		}
	}
	columnStr = strings.Join(resultFields, ", ")
	return
}

func buildWheres(conditions []Condition, wheres []Where) (whereStr string, err error) {
	var whereStrs []string

	wheresMap := map[string]Where{}
	for _, w := range wheres {
		wheresMap[w.Field] = w
	}

	conditionsMap := map[string]Condition{}
	for _, c := range conditions {
		conditionsMap[c.Name] = c
	}

	for _, c := range conditions {

		where, ok := wheresMap[c.Name]

		if !ok {
			err = fmt.Errorf("%s not in schema ", c.Name)
			return
		}

		handlerStr := where.Handler
		if len(handlerStr) > 0 {
			handler := FieldHandler{
				Field: where.Field,
			}

			handlerStructStrs := strings.Split(handlerStr, ";")
			handlerAndMethodStrs := strings.Split(handlerStructStrs[0], ".")

			if len(handlerStructStrs) >= 2 && len(handlerAndMethodStrs) == 2 {

				handler.HandlerName = handlerAndMethodStrs[0]
				handler.MethodName = handlerAndMethodStrs[1]
				handlerMethodArgs := []interface{}{}
				if len(handlerStructStrs[1:]) > 0 {
					for _, v := range handlerStructStrs[1:] {
						value, ok := conditionsMap[v]
						if !ok {
							err = fmt.Errorf("field:%s not found in conditions", v)
							return
						}
						handlerMethodArgs = append(handlerMethodArgs, value)
					}
				}

				handlerResult, handErr := doHandler(&handler)
				if handErr != nil {
					return
				}
				c.Value = handlerResult
			} else {
				err = fmt.Errorf("schema error handler not conform to the standards ,column is %s .", where.Field)
				return
			}
		}

		switch c.Type {
		case "eq":
			str := fmt.Sprintf("%s = %v", c.Name, c.Value)
			whereStrs = append(whereStrs, str)
		case "ne":
			str := fmt.Sprintf("%s != %v", c.Name, c.Value)
			whereStrs = append(whereStrs, str)
		case "lt":
			str := fmt.Sprintf("%s > %v", c.Name, c.Value)
			whereStrs = append(whereStrs, str)
		case "gt":
			str := fmt.Sprintf("%s < %v", c.Name, c.Value)
			whereStrs = append(whereStrs, str)
		case "like":
			str := fmt.Sprintf("%s like %v", c.Name, c.Value)
			whereStrs = append(whereStrs, str)
		}

	}

	whereStr = strings.Join(whereStrs, " AND ")
	return
}

func buildJoins(tableName string, joins []Join) (joinStr string, err error) {
	resultStrs := []string{}

	for _, join := range joins {
		// LEFT JOIN demo_other ON demo_other.demo_id = demo.id

		str := fmt.Sprintf("%s %s ON %s.%s = %s.%s", join.JoinType, join.JoinTableName, join.JoinTableName, join.JoinTableField, tableName, join.LinkField)

		resultStrs = append(resultStrs, str)

	}

	joinStr = strings.Join(resultStrs, " ")
	return
}

func buildPage(page int, size int) string {
	offset := (page - 1) * size
	return fmt.Sprintf("LIMIT %v OFFSET %v", size, offset)
}

func buildOrders(arr []string) string {
	return strings.Join(arr, ", ")
}

func doHandler(h *FieldHandler) (result interface{}, err error) {
	// 获取处理器工厂
	factory, ok := handler.GetHandlerFactory(h.HandlerName)
	if !ok {
		err = fmt.Errorf("handler:%s not found", h.HandlerName)
		return
	}

	// 创建处理器实例
	instance, err := factory()
	if err != nil {
		return
	}

	// 使用反射查找方法
	method := reflect.ValueOf(instance).MethodByName(h.MethodName)
	if !method.IsValid() {
		err = fmt.Errorf("Method not found: %s\n", h.MethodName)
	}

	argList := []reflect.Value{}

	for _, v := range h.Args {
		argList = append(argList, reflect.ValueOf(v))
	}

	// 调用方法
	results := method.Call(argList)
	result = results[0].Interface()

	return
}
