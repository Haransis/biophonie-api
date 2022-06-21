package user

type AddUser struct {
	Name string `json:"name" example:"bob" binding:"required,min=3,max=20"`
}

type User struct {
	Id        int    `db:"id" json:"userId" example:"123"`
	Name      string `db:"name" json:"name" example:"bob" binding:"required"`
	Token     string `db:"token" json:"token" example:"auinrsetanruistnstnaustie"`
	CreatedOn string `db:"created_on" json:"createdOn" example:"2022-05-26T11:17:35.079344Z"`
}
