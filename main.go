package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	name   string
	scores []int
}

func (student *Student) addScore(score int) {
	student.scores = append(student.scores, score)
}

func (student *Student) averageScore() float64 {
	sum := 0
	for _, score := range student.scores {
		sum += score
	}

	avg := float64(sum) / float64(len(student.scores))
	return avg
}

func main() {
	file, err := os.Open("students.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	students := make(map[string]*Student)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		name := parts[0]

		// Оценка - это второй элемент (конвертируем из строки в число)
		score, _ := strconv.Atoi(parts[1])

		// Если студент уже существует в мапе, добавляем ему оценку
		if student, exists := students[name]; exists {
			student.addScore(score)
		} else {
			// Если студента ещё нет, создаём нового и добавляем оценку
			students[name] = &Student{name: name, scores: []int{score}}
		}
	}

	// Преобразуем словарь в список для сортировки
	var studentList []*Student
	for _, student := range students {
		studentList = append(studentList, student)
	}

	// Сортируем студентов по имени
	sort.Slice(studentList, func(i, j int) bool {
		return studentList[i].name < studentList[j].name
	})

	// Выводим данные о каждом студенте
	for _, student := range studentList {
		fmt.Println(student.name)

		fmt.Print("Scores: ")
		for i, score := range student.scores {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(score)
		}
		fmt.Println()

		fmt.Printf("Average score: %.2f\n", student.averageScore())

		fmt.Println()
	}

}
