package api

type Request struct {
	Content  string `json:"content" binding:"required"`
	Password string `json:"password" binding:"required"`
}
