package huskapi

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/serials"
	"time"
)

// Sample api is a collection of books
type context struct {
	Books husk.Tabler
}

var (
	ctx context
	open bool
)

func CreateContext() {
	if !open {
		ctx = context{
			Books: husk.NewTable(Book{}, serials.GobSerial{}),
		}

		seed()
		open = true
	}
}

func Shutdown() {
	ctx = context{}
	open = false
}

func (c context) Save() error {
	return ctx.Books.Save()
}

func seed() {
	err := ctx.Books.Seed("db/books.seed.json")

	if err != nil {
		panic(err)
	}

	ctx.Books.Save()
}

type Book struct {
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	SubTitle    string    `json:"subtitle"`
	Author      string    `json:"author"`
	Published   time.Time `json:"published"`
	Publisher   string    `json:"publisher"`
	Pages       int       `json:"pages"`
	Description string    `json:"description"`
	Website     string    `json:"website"`
}

func (o Book) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}