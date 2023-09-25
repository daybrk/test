package controller

type User struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic"`
}

type DeleteUser struct {
	Id int `json:"id"`
}

type ModifyUser struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  *string  `json:"patronymic"`
	Age         int      `json:"age"`
	Gender      string   `json:"gender"`
	Nationality []string `json:"nationality"`
}

type FilterUser struct {
	Name        *string  `json:"name,omitempty"`
	Surname     *string  `json:"surname,omitempty"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Age         *int     `json:"age,omitempty"`
	Gender      *string  `json:"gender,omitempty"`
	Nationality []string `json:"nationality,omitempty"`
}
