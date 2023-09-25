package postgresdb

type EnrichmentUser struct {
	Id          int      `db:"id"`
	Name        string   `db:"name"`
	Surname     string   `db:"surname"`
	Patronymic  *string  `db:"patronymic"`
	Age         int      `db:"age"`
	Gender      string   `db:"gender"`
	Nationality []string `db:"nationality"`
}

type Filter struct {
	Name        *string
	Surname     *string
	Patronymic  *string
	Age         *int
	Gender      *string
	Nationality []string
}
