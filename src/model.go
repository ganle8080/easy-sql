package easysql

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

// OperateSchema 通用操作schema结构体
type OperateSchema struct{}
