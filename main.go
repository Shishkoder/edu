package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	name   string
	scores map[string][]int // Пример: {"math": [5, 2, 3]}
}

func (student *Student) addScore(discipline string, score int) {
	student.scores[discipline] = append(student.scores[discipline], score)
}

func (student *Student) averageScore() float64 {
	totalSum := 0.0
	count := 0
	for _, scores := range student.scores {
		// Если оценок нет по дисциплине, то пропускаем
		if len(scores) == 0 {
			continue
		}
		summa := 0
		for _, score := range scores {
			summa += score
		}
		avg := float64(summa) / float64(len(scores)) // Средний балл по дисциплине
		totalSum += avg
		count++
	}
	// Если дисциплин вообще нет
	if count == 0 {
		return 0
	}
	return totalSum / float64(count) // Средний балл по всем дисциплинам
}

// Функция для чтения данных из файла или консоли
func readStudentData(inputFile string) (map[string]*Student, error) {
	students := make(map[string]*Student)
	validDisciplines := map[string]bool{
		"math":    true,
		"english": true,
		"physics": true,
	}

	var scanner *bufio.Scanner
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
		fmt.Println("Чтение данных из файла:", inputFile)
	} else {
		fmt.Println("Введите данные о студентах построчно (имя, дисциплина, оценка). Пустая строка завершает ввод:")
		scanner = bufio.NewScanner(os.Stdin)
	}

	// Чтение данных
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			fmt.Println("Ошибка: строка должна содержать имя, дисциплину и оценку")
			continue
		}

		name := parts[0]
		discipline := parts[1]
		score, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("Ошибка: оценка должна быть числом")
			continue
		}

		if !validDisciplines[discipline] {
			fmt.Printf("Ошибка: дисциплина '%s' не поддерживается\n", discipline)
			continue
		}

		// Добавление студента или обновление существующего
		if student, exists := students[name]; exists {
			student.addScore(discipline, score)
		} else {
			students[name] = &Student{name: name, scores: make(map[string][]int)}
			students[name].addScore(discipline, score)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении данных: %v", err)
	}

	return students, nil
}

// Вывод данных о студентах (в файл или консоль)
func writeStudentData(students []*Student, outputFile string) error {
	var writer *bufio.Writer

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("ошибка при создании файла: %v", err)
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
		fmt.Println("Запись в файл:", outputFile)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}

	for _, student := range students {
		var builder strings.Builder

		builder.WriteString(student.name + "\n")
		builder.WriteString("Scores: \n")

		for discipline, scores := range student.scores {
			builder.WriteString(fmt.Sprintf("%s: ", discipline))

			for i, score := range scores {
				if i > 0 {
					builder.WriteString(", ")
				}
				builder.WriteString(fmt.Sprintf("%d", score))
			}
			builder.WriteString("\n")
		}

		builder.WriteString(fmt.Sprintf("Average score: %.2f\n", student.averageScore()))
		builder.WriteString("\n")

		writer.WriteString(builder.String())
	}

	return writer.Flush()
}

// Сортировка студентов по среднему баллу и имени
func sortStudents(students map[string]*Student) []*Student {
	var studentList []*Student
	for _, student := range students {
		studentList = append(studentList, student)
	}

	sort.Slice(studentList, func(i, j int) bool {
		// Если ср. балл разный, то выводим в порядке убывания
		if studentList[i].averageScore() != studentList[j].averageScore() {
			// Сортировка по среднему баллу в порядке убывания
			return studentList[i].averageScore() > studentList[j].averageScore()
		}
		// Если ср. балл одинаковый, сортируем по имени в порядке возрастания
		return studentList[i].name < studentList[j].name
	})

	return studentList
}

func main() {
	outputFlag := flag.String("o", "", "Output file")
	outputFlagLong := flag.String("output", "", "Output file")

	inputFlag := flag.String("f", "", "Input file")
	inputFlagLong := flag.String("file", "", "Input file")

	flag.Parse()

	// Определяем файл для ввода данных
	inputFile := ""
	if *inputFlag != "" {
		inputFile = *inputFlag
	} else if *inputFlagLong != "" {
		inputFile = *inputFlagLong
	}

	// Определяем файл для вывода данных
	outputFile := ""
	if *outputFlag != "" {
		outputFile = *outputFlag
	} else if *outputFlagLong != "" {
		outputFile = *outputFlagLong
	}

	// Чтение данных о студентах
	students, err := readStudentData(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Сортировка студентов
	sortedStudents := sortStudents(students)

	// Вывод данных о студентах
	err = writeStudentData(sortedStudents, outputFile)
	if err != nil {
		fmt.Println(err)
	}
}
