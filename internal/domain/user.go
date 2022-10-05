package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	VkId     string `json:"vk_id"`
	GoogleId string `json:"google_id"`
	GithubId string `json:"github_id"`
	AppleId  string `json:"apple_id"`
}