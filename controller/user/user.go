package user

type AddUser struct {
	Username string `json:"username" example:"bob" binding:"required"`
}

type User struct {
	UserId    *int    `db:"user_id" json:"user_id" example:"123"`
	UserName  *string `db:"username" json:"username" example:"bob" binding:"required"`
	Token     *string `db:"token" json:"token" example:"auinrsetanruistnstnaustie"`
	CreatedOn *string `db:"created_on" json:"created_on" example:"2022-05-26T11:17:35.079344Z"`
	LastLogin *string `db:"last_login" json:"last_login" example:"2022-05-26T11:17:35.079344Z"`
}
