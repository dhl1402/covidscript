package playground

type (
	InterpretRequest struct {
		Script string
	}
	InterpretResponse struct {
		Response []string
		Error    string
	}
)
