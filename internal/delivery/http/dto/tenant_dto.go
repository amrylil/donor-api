package dto

// DTO untuk request body (Create & Update)
type TenantRequest struct {
	Name string `json:"name" binding:"required"`
}

// DTO untuk response (data aman untuk publik)
type TenantResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
