package api

import (
	"FakeAPI/internal/logger"
	"FakeAPI/internal/service"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sync"
	"time"
)

func apiLogging(method, url string) func() {
	startTime := time.Now()
	return func() {
		duration := time.Since(startTime).Milliseconds()
		logger.GetLogger().Info(
			"API timing",
			zap.String("Method", method),
			zap.String("URL", url),
			zap.Int64("API duration", duration),
		)
	}
}

func GetAllSync(w http.ResponseWriter, r *http.Request) {
	endTiming := apiLogging("GET", "/sync")
	_, err := service.GetBooks()
	if err != nil {
		logger.GetLogger().Error("Unable to get data", zap.String("error", err.Error()))
	}
	_, err = service.GetUsers()
	if err != nil {
		logger.GetLogger().Error("Unable to get data", zap.String("error", err.Error()))
	}
	_, err = service.GetProducts()
	if err != nil {
		logger.GetLogger().Error("Unable to get data", zap.String("error", err.Error()))
	}
	endTiming()
}

//type Response struct {
//	name string
//	data *http.Response
//}

//func GetAllAsync(w http.ResponseWriter, r *http.Request) {
//	endTiming := apiLogging("GET", "/async")
//	ch := make(chan *Response)
//	reqMap := map[string]func() (res *http.Response, err error){
//		"users":    service.GetUsers,
//		"products": service.GetProducts,
//		"books":    service.GetBooks,
//	}
//	for resource, _ := range reqMap {
//		go func(resource string) {
//			res, _ := reqMap[resource]()
//			resp := &Response{name: resource, data: res}
//			ch <- resp
//		}(resource)
//	}
//	data := make(map[string]interface{})
//	for i := 0; i < 3; i++ {
//		select {
//		case res := <-ch:
//			func() {
//				var resData interface{}
//				defer res.data.Body.Close()
//				if err := json.NewDecoder(res.data.Body).Decode(&resData); err != nil {
//					data[res.name] = nil
//				} else {
//					data[res.name] = resData
//				}
//			}()
//		}
//	}
//	res, _ := json.Marshal(data)
//	w.Write(res)
//	endTiming()
//}

func GetAllAsync(w http.ResponseWriter, r *http.Request) {
	endTiming := apiLogging("GET", "/async")
	var tasks sync.WaitGroup
	var results sync.Map
	reqMap := map[string]func() (res *http.Response, err error){
		"users":    service.GetUsers,
		"products": service.GetProducts,
		"books":    service.GetBooks,
	}
	for resource, requestFunc := range reqMap {
		tasks.Add(1)
		go func(resource string, requestFunc func() (res *http.Response, err error)) {
			defer tasks.Done()
			res, err := requestFunc()
			if err != nil {
				results.Store(resource, nil)
			}
			if res == nil {
				log.Printf("Received nil response for %s", resource)
				results.Store(resource, nil)
				return
			}
			defer res.Body.Close()
			var resData interface{}
			err = json.NewDecoder(res.Body).Decode(&resData)
			if err != nil {
				results.Store(resource, nil)
			} else {
				results.Store(resource, resData)
			}
		}(resource, requestFunc)
	}
	tasks.Wait()
	data := make(map[string]interface{})
	for resource := range reqMap {
		apiRes, _ := results.Load(resource)
		if apiRes == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data[resource] = apiRes
	}
	res, _ := json.Marshal(data)
	w.Write(res)
	endTiming()
}

type ProductComparisonResponse struct {
	Price ProductPriceComparisonResponse `json:"price"`
}

type ProductPriceComparisonResponse struct {
	ProductOne float32 `json:"productOne"`
	ProductTwo float32 `json:"productTwo"`
	Winner     int     `json:"winner"`
}

func CompareProducts(w http.ResponseWriter, r *http.Request) {
	endTiming := apiLogging("GET", "/products/compare")
	productId1 := r.URL.Query().Get("product1")
	productId2 := r.URL.Query().Get("product2")
	if productId1 == "" || productId2 == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	productA, err := service.GetProductById(productId1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.GetLogger().Error("error while getting product with id: "+productId1, zap.Error(err))
		return
	}
	productB, err := service.GetProductById(productId2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.GetLogger().Error("error while getting product with id: "+productId2, zap.Error(err))
		return
	}
	var productCompResp ProductComparisonResponse
	var productPriceCompResp ProductPriceComparisonResponse
	if productA.Price <= productB.Price {
		productPriceCompResp.ProductOne = productA.Price
		productPriceCompResp.ProductTwo = productB.Price
		productPriceCompResp.Winner = productA.Id
	} else {
		productPriceCompResp.ProductOne = productA.Price
		productPriceCompResp.ProductTwo = productB.Price
		productPriceCompResp.Winner = productB.Id
	}
	productCompResp.Price = productPriceCompResp
	res, _ := json.Marshal(productCompResp)
	w.Write(res)
	endTiming()
}
