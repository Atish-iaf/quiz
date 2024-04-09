package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultTimeLimit = 10
	defaultFilePath  = "questionAnswers.csv"
)

type questionAnswer struct {
	question string
	answer   string
}

func main() {
	filePath := flag.String("file", defaultFilePath, "path of file which contains questions and answers.")
	timeLimit := flag.Int("timeLimit", defaultTimeLimit, "optional, time limit in seconds for the quiz. Default is 10s")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s : %s", *filePath, err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read from file: %s : %s", *filePath, err)
	}

	questionAnswers := getQuestionAnswers(lines)
	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))

	timeUp, correctAnswerCount := play(questionAnswers, timer)
	if timeUp {
		fmt.Println("Time up!")
	}
	fmt.Printf("You scored %d out of %d\n", correctAnswerCount, len(questionAnswers))
}

func getQuestionAnswers(lines [][]string) []questionAnswer {
	questions := make([]questionAnswer, 0)
	for _, line := range lines {
		questions = append(questions, questionAnswer{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		})
	}
	return questions
}

func play(questionAnswers []questionAnswer, timer *time.Timer) (timeUp bool, correct int) {
	correctAnswerCount := 0
	answerCh := make(chan string)
	for quesNum, questionAnswer := range questionAnswers {
		go getAnswer(quesNum+1, questionAnswer.question, answerCh)
		select {
		case <-timer.C:
			return true, correctAnswerCount
		case userAnswer := <-answerCh:
			if userAnswer == questionAnswer.answer {
				correctAnswerCount++
			}
		}
	}
	return false, correctAnswerCount
}

func getAnswer(quesNum int, question string, answerCh chan string) {
	fmt.Printf("%d. %s \n", quesNum, question)
	var userAnswer string
	fmt.Scanf("%s", &userAnswer)
	answerCh <- userAnswer
}
