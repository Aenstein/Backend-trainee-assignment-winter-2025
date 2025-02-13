package database

type Shop struct {
	Name  string
	Price int
}

type User struct {
	Username    string
	Coins       int
	Transaction []Transactions
}

type Transactions struct {
}
