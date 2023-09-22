package controller

type Fio struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type DeleteFio struct {
	Id int `json:"id"`
}

type ModifyFio struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Patronymic  string   `json:"patronymic"`
	Age         int      `json:"age"`
	Gender      string   `json:"gender"`
	Nationality []string `json:"nationality"`
}
