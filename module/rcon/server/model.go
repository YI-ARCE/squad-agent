package server

// LoginRequest 客户端请求登录
type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// LoginResponse 响应客户端
type LoginResponse struct {
	LoginStatus bool              `json:"login_status"`
	ExpireTime  int64             `json:"expire_time"`
	AId         string            `json:"a_id"`
	Role        map[string]string `json:"role"`
	Msg         string            `json:"msg"`
}

type Response struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Tag  string      `json:"tag,omitempty"`
}

type ResponseTag struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`
}
