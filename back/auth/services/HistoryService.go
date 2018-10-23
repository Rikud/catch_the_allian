package services

import "IT-Berries_Go_server/auth/DA"

func GetBestScoreForUserById(userId int) int {
	return DA.GetBestScoreForUserById(userId)
}
