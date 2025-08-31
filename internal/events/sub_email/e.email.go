package sub_email

type SendEvent struct {
	Email       string      `json:"email"`
	Template    string      `json:"template"`
	Locale      string      `json:"locale"`
	Data        interface{} `json:"data"`
	OperationId string      `json:"operation_id"`
}
