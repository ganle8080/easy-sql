package easysql

import (
	"strings"

	"gorm.io/gorm"
)

func Search(db *gorm.DB, searchData *SearchData, schema *SearchSchema) {

}

func dealWhere(searchFields []string, columns []Column) (resultFields []string, handlers []FieldHandler) {
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
			}
		}
	}
	return
}

func dealJoin() {

}
