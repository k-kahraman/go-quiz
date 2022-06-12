package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

func (p *Problem) ValidateAnswer(answer string) bool {
	return strings.ToLower(strings.Trim(answer, "\n\r\t")) == strings.ToLower(strings.Trim(p.Answer, "\n\r\t"))
}

type Quiz struct {
	Problems []Problem
	Score    uint
}

func (quiz *Quiz) IncrementScore(flags QuizFlags) {
	quiz.Score += *flags.CorrectScore
}

func (quiz *Quiz) PromptScoreMessage(scorePerCorrect int) {
	fmt.Printf("\nYou scored %d out of %d!", quiz.Score, len(quiz.Problems)*scorePerCorrect)
}

type QuizFlags struct {
	TimeLimit    *int
	FileName     *string
	CorrectScore *uint
}

func main() {
	quizFlags := InitFlags()
	quiz := InitQuiz(quizFlags)

	timer := time.NewTimer(time.Second * time.Duration(*quizFlags.TimeLimit))
	defer timer.Stop()

	go func() {
		<-timer.C
		quiz.PromptScoreMessage(int(*quizFlags.CorrectScore))
		ExitOnTimeUp()
	}()

	PlayQuiz(&quiz, &quizFlags)
	quiz.PromptScoreMessage(int(*quizFlags.CorrectScore))
	ExitOnTimeUp()
}

func InitQuiz(flags QuizFlags) Quiz {
	csvData := ReadCsvData(flags)
	return Quiz{
		Problems: ReadProblemsFromCsv(csvData),
		Score:    0,
	}
}

func InitFlags() QuizFlags {
	defer flag.Parse()
	return QuizFlags{
		TimeLimit:    flag.Int("time", 30, "Sets the time limit to given amount in seconds."),
		FileName:     flag.String("filename", "./problems.csv", "Path for the Problems CSV file."),
		CorrectScore: flag.Uint("score", 5, "Score for per correct answer."),
	}
}

func PlayQuiz(quiz *Quiz, flags *QuizFlags) {
	reader := bufio.NewReader(os.Stdin)

	for index, problem := range quiz.Problems {

		fmt.Printf("Q%d: %s? ", index+1, problem.Question)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if problem.ValidateAnswer(answer) {
			fmt.Println("Correct!")
			quiz.IncrementScore(*flags)
		} else {
			fmt.Printf("Incorrect answer! Correct answer was %s.\n", problem.Answer)
		}
	}
}

func PanicInErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadCsvData(flags QuizFlags) [][]string {
	fileName := *flags.FileName

	ioInput, ioErr := ioutil.ReadFile(fileName)
	PanicInErr(ioErr)

	csvInput := csv.NewReader(strings.NewReader(string(ioInput)))

	csvData, csvErr := csvInput.ReadAll()
	PanicInErr(csvErr)

	return csvData
}

func ReadProblemsFromCsv(csvData [][]string) []Problem {
	problems := make([]Problem, len(csvData))

	for index, data := range csvData {
		problems[index] = Problem{
			Question: data[0],
			Answer:   data[1],
		}
	}

	return problems
}

func ExitOnTimeUp() {
	fmt.Println("\nTime is up!")
	os.Exit(0)
}
