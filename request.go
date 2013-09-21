package barbershop

type Request struct {
	answer   chan string
	question string
}

func (self Request) GetAnswerChannel() {
	return answer
}

func (self Request) getQuestion() {
	return question
}
