package objects

// BaseBoolInt type.
type BaseBoolInt bool

// UsersOccupation struct.
type UsersOccupation struct {
	// BUG(VK): UsersOccupation.ID is float https://vk.com/bug136108
	ID   float64 `json:"id"`   // ID of school, university, company group
	Name string  `json:"name"` // Name of occupation
	Type string  `json:"type"` // Type of occupation
}

// UsersUser struct.
//schemaverify:schema_name=users_user_full
type UsersUser struct { // want `id field is required in object UsersUser \(users_user_full\)`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Domain      string          `json:"domain"`
	ScreenName  string          `json:"screen_name"`
	Bdate       string          `json:"bdate"`
	Photo50     string          `json:"photo_50"`
	Photo100    string          `json:"photo_100"`
	Online      BaseBoolInt     `json:"online"`
	Blacklisted BaseBoolInt     `json:"blacklisted"`
	Occupation  UsersOccupation `json:"occupation"`
}
