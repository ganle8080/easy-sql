package easysql

import (
	"encoding/json"
	"os"
)

// SearchSchema 通用查询schema结构体
type SearchSchema struct {
	TableName  string      `json:"table_name"`
	Columns    []Column    `json:"columns"`
	Joins      []Join      `json:"joins"`
	Conditions []Condition `json:"conditions"`
	Orders     []string    `json:"orders"`
	Groups     []string    `json:"groups"`
}

// Column 字段
type Column struct {
	Field       string `json:"field"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Handler     string `json:"handler"`
}

// Condition 条件
type Condition struct {
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

// OperateSchema 通用操作schema结构体
type OperateSchema struct{}

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
