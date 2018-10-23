package services

import (
	"IT-Berries_Go_server/auth/DA"
	"IT-Berries_Go_server/auth/models"
)

func FindUserByEmail(email string) *models.User {
	searchResult := DA.FindUserByEmail(email)
	if len(searchResult) > 1 || len(searchResult) < 1 {
		return nil
	}
	return searchResult[0]
}

func AddUser(user models.User) {
	result := DA.AddUser(user)
	if result < 0 {
		panic("Error trying to add user data to database!")
	}
}

func findUserById(userId int) *models.User {
	searchResult := DA.FindUserById(userId)
	if len(searchResult) > 1 || len(searchResult) < 1 {
		return nil
	}
	return searchResult[0]
}

func SaveUser(user models.User) {
	result := DA.SaveUser(user)
	if result < 0 {
		panic("Error trying to save user data to database!")
	}
}

func FindAllUsersForScoreBoard() []*models.ScoreRecord {
	searchResult := DA.GetScoreBoardData()
	if len(searchResult) < 0 {
		panic("Error trying to save user data to database!")
	}
	return searchResult
}