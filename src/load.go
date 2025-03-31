package easysql

import (
	"encoding/json"
	"os"
)

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
