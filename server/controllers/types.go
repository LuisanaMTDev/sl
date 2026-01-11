package controllers

type AddLessonRequest struct {
	Topic            string   `json:"topic"`
	AmountOfCards    int64    `json:"amount_of_cards"`
	RepetitionsDates []string `json:"repetitions_dates"`
}

type UserLoginInfoRequest struct {
	Password string `json:"password"`
}
