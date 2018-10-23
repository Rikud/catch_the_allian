package models

type ScoreRecord struct {
	Id int			`json:"id"`
	Username string `json:"username"`
	Score int		`json:"score"`
}


func (record *ScoreRecord) GetUsernamePoint() *string {
	return &record.Username
}

func (record *ScoreRecord) GetScorePoint() *int {
	return &record.Score
}

func (record *ScoreRecord) GetUsername() string {
	return record.Username
}

func (record *ScoreRecord) GetScore() int {
	return record.Score
}

func (record *ScoreRecord) SetUsername(username string) {
	record.Username = username
}

func (record *ScoreRecord) SetScore(score int) {
	record.Score = score
}

func (record *ScoreRecord) SetId(id int) {
	record.Id = id
}