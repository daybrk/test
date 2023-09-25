package user

type User struct {
	Name       string
	Surname    string
	Patronymic *string
}

type EnrichmentUser struct {
	Id          int
	Name        string
	Surname     string
	Patronymic  *string
	Age         int
	Gender      string
	Nationality []string
}

type Filter struct {
	Name        *string
	Surname     *string
	Patronymic  *string
	Age         *int
	Gender      *string
	Nationality []string
}
