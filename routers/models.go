package routers

type SignUpPayload struct {
	Instance string `json:"instance"`
	Token    string `json:"token"`
}
