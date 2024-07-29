package service

import (
	"FakeAPI/internal/cache"
	"FakeAPI/internal/network"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var (
	BASE_URL = "https://fakestoreapi.com"
	client   = network.NewClient().Client("fakestore").Timeout(10)
)

type Product struct {
	Id    int     `json:"id"`
	Title string  `json:"title"`
	Price float32 `json:"price"`
}

func isNum(a string) bool {
	_, err := strconv.Atoi(a)
	return err == nil
}

var lruCache *cache.LRUCache

func init() {
	if lruCache == nil {
		lruCache = cache.NewLRUCache(2)
	}
}

func getProductFromResponse(res *http.Response) (*Product, error) {
	defer res.Body.Close()
	var product Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func GetProductById(id string) (*Product, error) {
	if !isNum(id) {
		return nil, errors.New("id must be an integer")
	}
	idInt, _ := strconv.Atoi(id)
	cachedPrice, err := lruCache.Get(idInt)
	if err != nil {
		res, err := client.Get(BASE_URL + "/products/" + id)
		if err != nil {
			return nil, err
		}
		if res == nil || (res.StatusCode < 200 || res.StatusCode > 299) {
			return nil, err
		}
		product, productErr := getProductFromResponse(res)
		lruCache.Set(idInt, product.Price)
		return product, productErr
	}
	return &Product{Id: idInt, Price: cachedPrice}, nil
}
