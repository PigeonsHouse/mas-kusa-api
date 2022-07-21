package utils

type SignUpPayload struct {
	Instance string `json:"instance"`
	Token    string `json:"token"`
}

type JwtInfo struct {
	Jwt string `json:"jwt"`
}

type ImagePath struct {
	Path    string `json:"path"`
	Refresh bool   `json:"refresh"`
}
