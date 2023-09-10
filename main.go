package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"io/ioutil"
	"io"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)
var fil string
func main() {
	a := app.New()
	w := a.NewWindow("File Processing App")
	w.Resize(fyne.NewSize(800, 500))

	entry := widget.NewMultiLineEntry()
	entry.Resize(fyne.NewSize(600, 300))
	entry.Move(fyne.NewPos(100, 135))

	var fileContent string   // Переменная для хранения содержимого файла
	var selectedAction string // Переменная для хранения выбранного действия

	// Создаем выпадающий список для выбора действия
	actionSelect := widget.NewSelect([]string{
		"Группировка по расстоянию",
		"Группировка по имени",
		"Группировка по времени создания",
		"Группировка по типу",
	}, func(selected string) {
		// При выборе действия обновляем переменную selectedAction
		selectedAction = selected
	})

	actionSelect.Resize(fyne.NewSize(200, 30))
	actionSelect.Move(fyne.NewPos(100, 30))

	btn := widget.NewButton("Open File", func() {
		dialog.ShowFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if err == nil && r != nil {
					data, _ := ioutil.ReadAll(r.(io.Reader))
	
					fileContent = string(data) // Сохраняем содержимое файла в переменной
					entry.SetText(fileContent) // Отображаем содержимое в многострочном поле
	
					// Сохраняем путь к файлу в переменной fil
					fil = r.URI().String()
	
					r.Close()
				}
			},
			w,
		)
	})
	btn.Resize(fyne.NewSize(150, 75))
	btn.Move(fyne.NewPos(325, 30))

	applyBtn := widget.NewButton("Применить", func() {
		if selectedAction != "" && fileContent != "" {
			processFileContent(selectedAction, fileContent)
			dialog.ShowInformation("Успешно выполнено", "Действие выполнено успешно и сохранено в output.txt", w)
		} else {
			err := fmt.Errorf("Выберите действие и откройте файл")
			dialog.ShowError(err, w)
		}
	})
	applyBtn.Resize(fyne.NewSize(150, 75))
	applyBtn.Move(fyne.NewPos(500, 30))

	cont := container.NewWithoutLayout(
		actionSelect,
		btn,
		applyBtn,
		entry,
	)

	w.SetContent(cont)
	w.ShowAndRun()
}

func processFileContent(selectedAction, fileContent string) {
	// Выполните обработку содержимого файла в зависимости от выбранного действия
	switch selectedAction {
	case "Группировка по расстоянию":
		objects := parseFileContent(fileContent) // Парсим содержимое файла в объекты
		SortAndWriteByDistance(objects, "output.txt")      // Выполняем группировку по расстоянию
	case "Группировка по имени":
		objects := parseFileContent(fileContent) // Парсим содержимое файла в объекты
		grname(objects, "output.txt")                     // Выполняем группировку по имени
	case "Группировка по времени создания":
		objects := parseFileContent(fileContent) // Парсим содержимое файла в объекты
		grtime(objects, "output.txt")                     // Выполняем группировку по времени создания
	case "Группировка по типу":
		objects := parseFileContent(fileContent) // Парсим содержимое файла в объекты
		grtype(objects, "output.txt")                     // Выполняем группировку по типу
	}
}


func parseFileContent(fileContent string) []map[string]string {
	lines := strings.Split(fileContent, "\n")
	var objects []map[string]string

	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) >= 5 {
			obj := map[string]string{
				"name": fields[0],
				"x":    fields[1],
				"y":    fields[2],
				"typ":  fields[3],
				"time": fields[4],
			}

			objects = append(objects, obj)
		}
	}

	return objects
}

func SortAndWriteByDistance(objects []map[string]string, fil string) {
	// Функция для вычисления расстояния от (0,0) до точки (x, y)
	distance := func(x, y string) float64 {
		xFloat, errX := strconv.ParseFloat(x, 64)
		yFloat, errY := strconv.ParseFloat(y, 64)
		if errX != nil || errY != nil {
			// Обработка ошибок преобразования
			fmt.Printf("Ошибка преобразования координат для объекта\n")
			return math.MaxFloat64 // Возвращаем максимальное значение, чтобы объект считался слишком далеко
		}
		return math.Sqrt(xFloat*xFloat + yFloat*yFloat)
	}

	// Списки для разделения данных
	var upTo100, upTo1000, upTo10000 []map[string]string

	// Разделяем объекты по расстояниям
	for _, obj := range objects {
		dist := distance(obj["x"], obj["y"])
		switch {
		case dist < 100:
			upTo100 = append(upTo100, obj)
		case dist < 1000:
			upTo1000 = append(upTo1000, obj)
		case dist < 10000:
			upTo10000 = append(upTo10000, obj)
		default:
			upTo10000 = append(upTo10000, obj)
		}
	}

	// Сортировка каждой группы по расстоянию
	sort.Slice(upTo100, func(i, j int) bool {
		dist1 := distance(upTo100[i]["x"], upTo100[i]["y"])
		dist2 := distance(upTo100[j]["x"], upTo100[j]["y"])
		return dist1 < dist2
	})

	sort.Slice(upTo1000, func(i, j int) bool {
		dist1 := distance(upTo1000[i]["x"], upTo1000[i]["y"])
		dist2 := distance(upTo1000[j]["x"], upTo1000[j]["y"])
		return dist1 < dist2
	})

	sort.Slice(upTo10000, func(i, j int) bool {
		dist1 := distance(upTo10000[i]["x"], upTo10000[i]["y"])
		dist2 := distance(upTo10000[j]["x"], upTo10000[j]["y"])
		return dist1 < dist2
	})

	// Создаем или открываем файл для записи
	file, err := os.Create(fil)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	// Записываем данные в файл с разделами
	writeGroup := func(group []map[string]string, label string) {
		if len(group) > 0 {
			fmt.Fprintf(file, "=== Расстояние %s ===\n", label)
			for _, obj := range group {
				fmt.Fprintf(file, "%s, Расстояние: %.2f, Тип: %s, Время: %s\n", obj["name"], distance(obj["x"], obj["y"]), obj["typ"], obj["time"])
			}
		}
	}

	writeGroup(upTo100, "до 100")
	writeGroup(upTo1000, "до 1000")
	writeGroup(upTo10000, "до 10000")

	fmt.Println("Результаты записаны в файл.")
}

func grname(objects []map[string]string, fil string) {
	// Создаем карту, где ключ - это имя группы, значение - это список объектов
	groups := make(map[string][]map[string]string)

	// Группировка объектов по именам
	for _, obj := range objects {
		name := obj["name"]
		groupName := ""

		// Проверяем, что все символы имени являются кириллическими
		isCyrillic := true
		for _, char := range name {
			if !unicode.Is(unicode.Cyrillic, char) {
				isCyrillic = false
				break
			}
		}

		if isCyrillic {
			groupName = "" // Сбрасываем название группы для кириллических имен
		} else if groupName == "" {
			groupName = "#" // Устанавливаем группу "#" только один раз
		}

		groups[groupName] = append(groups[groupName], obj)
	}

	// Создаем или открываем файл для записи
	file, err := os.Create(fil)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	// Записываем данные в файл
	for groupName, group := range groups {
		if groupName != "" {
			fmt.Fprintf(file, "Группа %s:\n", groupName)
		}
		for _, obj := range group {
			fmt.Fprintf(file, "%s, Расстояние: %s, Тип: %s, Время: %s\n", obj["name"], obj["distance"], obj["typ"], obj["time"])
		}
	}

	fmt.Println("Результаты группировки записаны в файл.")
}

func grtime(objects []map[string]string, fil string) {
	// Создаем карту, где ключ - это название группы времени, значение - это список объектов
	groups := make(map[string][]map[string]string)

	// Группировка объектов по времени
	for _, obj := range objects {
		createTimeUnix, err := strconv.ParseFloat(obj["time"], 64)
		if err != nil {
			fmt.Printf("Ошибка парсинга времени для объекта %s: %v\n", obj["name"], err)
			continue
		}

		createTime := time.Unix(int64(createTimeUnix), 0)
		createTimeFormatted := createTime.Format("2006-01-02 15:04:05")

		var groupName string

		// Определяем группу времени в зависимости от временной разницы
		if createTime.After(time.Now().Add(-24 * time.Hour)) {
			groupName = "Сегодня"
		} else if createTime.After(time.Now().Add(-48 * time.Hour)) {
			groupName = "Вчера"
		} else if createTime.After(time.Now().Add(-7 * 24 * time.Hour)) {
			groupName = "На этой неделе"
		} else if createTime.After(time.Now().Add(-30 * 24 * time.Hour)) {
			groupName = "В этом месяце"
		} else if createTime.After(time.Now().Add(-365 * 24 * time.Hour)) {
			groupName = "В этом году"
		} else {
			groupName = "Ранее"
		}

		objWithFormattedTime := make(map[string]string)
		objWithFormattedTime["name"] = obj["name"]
		objWithFormattedTime["distance"] = obj["distance"]
		objWithFormattedTime["typ"] = obj["typ"]
		objWithFormattedTime["time"] = createTimeFormatted

		groups[groupName] = append(groups[groupName], objWithFormattedTime)
	}

	// Создаем или открываем файл для записи
	file, err := os.Create(fil)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	// Записываем данные в файл
	for groupName, group := range groups {
		if groupName != "" {
			fmt.Fprintf(file, "Группа %s:\n", groupName)
		}
		for _, obj := range group {
			fmt.Fprintf(file, "%s, Расстояние: %s, Тип: %s, Время: %s\n", obj["name"], obj["distance"], obj["typ"], obj["time"])
		}
	}

	fmt.Println("Результаты группировки по времени записаны в файл.")
}

func grtype(objects []map[string]string, fil string) {
	var minCount int
	fmt.Scan(&minCount)
	// Создаем карту, где ключ - это тип объекта, значение - это список объектов этого типа
	groups := make(map[string][]map[string]string)

	// Группировка объектов по типу
	for _, obj := range objects {
		objType := obj["typ"]

		// Добавляем объект в соответствующую группу по типу
		groups[objType] = append(groups[objType], obj)
	}

	// Удаляем группы, в которых меньше объектов, чем minCount
	for typ, group := range groups {
		if len(group) < minCount {
			delete(groups, typ)
		}
	}

	// Сортируем объекты внутри каждой группы по имени
	for _, group := range groups {
		sort.Slice(group, func(i, j int) bool {
			return group[i]["name"] < group[j]["name"]
		})
	}

	// Создаем или открываем файл для записи
	file, err := os.Create(fil)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	// Записываем данные в файл
	for typ, group := range groups {
		fmt.Fprintf(file, "Группа %s:\n", typ)
		for _, obj := range group {
			fmt.Fprintf(file, "%s, Расстояние: %s, Тип: %s, Время: %s\n", obj["name"], obj["distance"], obj["typ"], obj["time"])
		}
	}

	fmt.Println("Результаты группировки по типу записаны в файл.")
}
