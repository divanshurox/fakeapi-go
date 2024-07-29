package api

import (
	"FakeAPI/internal/db"
	"FakeAPI/internal/logger"
	"FakeAPI/internal/mongo"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type AddUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UpsertUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func AddUserToDB(w http.ResponseWriter, r *http.Request) {
	var request AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.GetLogger().Error("Unable to parse request body", zap.Error(err))
		return
	}
	user := &User{
		Name: request.Name,
		Age:  request.Age,
	}
	database := db.GetDatabase(mongo.GetInstance())
	query := db.NewQuery().WithDatabase("course").WithCollection("users").WithObject(user)
	err := database.Insert(query)
	if err != nil {
		logger.GetLogger().Error("Error inserting to database", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpsertUserToDB(w http.ResponseWriter, r *http.Request) {
	var request UpsertUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.GetLogger().Error("Unable to parse request body", zap.Error(err))
		return
	}
	database := db.GetDatabase(mongo.GetInstance())
	query := db.NewQuery().WithDatabase("course").WithCollection("users").WithFilter(map[string]interface{}{
		"name": request.Name,
	}).WithObject(&User{
		Name: request.Name,
		Age:  request.Age,
	})
	err := database.Update(query)
	if err != nil {
		logger.GetLogger().Error("Error upserting to database", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUserFromDB(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	database := db.GetDatabase(mongo.GetInstance())
	query := db.NewQuery().WithDatabase("course").WithCollection("users").WithFilter(map[string]interface{}{
		"name": name,
	})
	var user User
	err := database.Get(query, &user)
	if err != nil || user == (User{}) {
		logger.GetLogger().Error("User not found", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.GetLogger().Error("Error encoding user response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
