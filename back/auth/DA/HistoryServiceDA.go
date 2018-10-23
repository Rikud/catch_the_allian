package DA

import "IT-Berries_Go_server/auth/models"

func GetBestScoreForUserById(userId int) int {
	db := connect()
	defer db.Close()
	rows, err := db.Query("select max(h.score) as score from History h , Users u where u.user_id = h.user_id \n"+
		"and  u.user_id = $1 and score > 0  group by u.user_name, score order by score desc", userId)
	errorCheck(err, executeError)
	defer rows.Close()
	scores := make([]int, 0)
	for rows.Next() {
		var score int
		err := rows.Scan(score)
		errorCheck(err, readRowError)
		scores = append(scores, score)
	}
	errorCheck(rows.Err(), readRowError)
	if len(scores) > 0 {
		return scores[0]
	} else {
		return 0
	}
}

func GetScoreBoardData() []*models.ScoreRecord {
	db := connect()
	defer db.Close()
	rows, err := db.Query("select u.user_name, max(h.score) as score "+
		"from History h , Users u where u.user_id = h.user_id group by u.user_name order by score desc")
	errorCheck(err, executeError)
	defer rows.Close()
	scoreRecords := make([]*models.ScoreRecord, 0)
	for i := 0; rows.Next(); i++ {
		scoreRecord:= new(models.ScoreRecord)
		err := rows.Scan(scoreRecord.GetUsernamePoint(), scoreRecord.GetScorePoint())
		scoreRecord.SetId(i + 1)
		errorCheck(err, readRowError)
		scoreRecords = append(scoreRecords, scoreRecord)
	}
	errorCheck(rows.Err(), readRowError)
	return scoreRecords
}