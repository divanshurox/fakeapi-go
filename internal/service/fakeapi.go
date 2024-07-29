package service

import (
	"FakeAPI/internal/network"
	"net/http"
)

var fakeAPIClient = network.NewClient().Client("fakerapi").
	Timeout(20)

func GetUsers() (res *http.Response, err error) {
	res, err = fakeAPIClient.
		Get("https://fakerapi.it/api/v1/persons")
	return
}

func GetProducts() (res *http.Response, err error) {
	res, err = fakeAPIClient.
		Get("https://fakerapi.it/api/v1/products")
	return
}

func GetBooks() (res *http.Response, err error) {
	res, err = fakeAPIClient.
		Get("https://fakerapi.it/api/v1/books")
	return
}
