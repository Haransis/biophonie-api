package user

type AddUser struct {
	Name string `json:"name" example:"bob" binding:"required"`
}

type User struct {
	Id        *int    `db:"id" json:"user_id" example:"123"`
	Name      *string `db:"name" json:"name" example:"bob" binding:"required"`
	Token     *string `db:"token" json:"token" example:"auinrsetanruistnstnaustie"`
	CreatedOn *string `db:"created_on" json:"created_on" example:"2022-05-26T11:17:35.079344Z"`
	LastLogin *string `db:"last_login" json:"last_login" example:"2022-05-26T11:17:35.079344Z"`
}
