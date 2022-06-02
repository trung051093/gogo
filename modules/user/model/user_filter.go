package usermodel

type UserFilter struct {
	Fields    []string `json:"fields"`
	SortField string   `json:"sortField"`
	SortName  string   `json:"sortName"`
}

func (u *UserFilter) Process() error {
	if len(u.Fields) == 0 {
		u.Fields = []string{
			UserField["id"],
			UserField["firstName"],
			UserField["lastName"],
			UserField["email"],
			UserField["address"],
			UserField["company"],
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

	if u.SortField == "" {
		u.SortField = "id"
	}

	if u.SortName == "" {
		u.SortName = "desc"
	}

	return nil
}
