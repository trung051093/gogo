package usermodel

type UserFilter struct {
	Fields []string `json:"fields"`
	Order  string   `json:"order"`
}

func (u *UserFilter) Process() error {
	if len(u.Fields) == 0 {
		u.Fields = []string{
			UserField["firstName"],
			UserField["lastName"],
			UserField["email"],
			UserField["gender"],
			UserField["phoneNumber"],
		}
	} else {
		for index, field := range u.Fields {
			if _, ok := UserField[field]; ok {
				u.Fields[index] = UserField[field]
			}
		}
	}

	if u.Order == "" {
		u.Order = "Id desc"
	}

	return nil
}
