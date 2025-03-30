package easysql

type SearchData struct {
	SchemaName string `json:"schema_name"`
	Fields     []string
	Conditions string
}
