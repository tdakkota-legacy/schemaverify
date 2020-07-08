package objects

type BaseBoolInt bool

type UsersOccupation struct {
	// Fields with "BUG*" are ignored
	// BUG(VK): UsersOccupation.ID is float https://vk.com/bug136108
	ID   float64 `json:"id"`
	Name string  `json:"name"`
	Type string  `json:"type"`
}

type usersUserPhoto struct {
	Photo50  int    `json:"photo_50"`
	Photo100 string `json:"photo_100"`
}

// Mapped to 'users_user_full', not to 'users_user'
//schemaverify:schema_name=users_user_full
type UsersUser struct { // want `'id' field is required in object 'UsersUser' \('users_user_full'\)`
	FirstName  string `json:"first_name"`
	LastName   int    `json:"last_name"` // want `Field 'LastName' does not match schema: expected string instead of int`
	Domain     string `json:"domain"`
	ScreenName string `json:"screen_name"`
	Bdate      string `json:"bdate"`
	// embedded struct
	usersUserPhoto
	Online      BaseBoolInt     `json:"online"`
	Blacklisted BaseBoolInt     `json:"blacklisted"`
	Occupation  UsersOccupation `json:"occupation"`
}

type notPublicStructWhichShouldNotBeChecked struct {
	Value string
}
