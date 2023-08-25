package rabbitmq

type Message struct {
	Type       string      `json:"type"`
	SubType    string      `json:"sub_type"`
	StructName string      `json:"struct_name"`
	Key        string      `json:"key"`
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
}
