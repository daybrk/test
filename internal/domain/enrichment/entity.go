package enrichment

type Fio struct {
	Name       string
	Surname    string
	Patronymic string
}

type FioEnrichment struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality []string
}
