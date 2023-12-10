package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"strconv"
	"strings"
)

// Employee структура для хранения данных о сотруднике
type Employee struct {
	Name     string
	Children []Child
}

// Child структура для хранения данных о ребенке
type Child struct {
	Name          string
	Age           int
	Appearance    string
	FavoriteColor string
	Comments      string
}

var employees []Employee
var employeeList *widget.List
var content *fyne.Container

func main() {
	loadData()

	myApp := app.New()
	myWindow := myApp.NewWindow("Список детей сотрудников авиакомпаний")

	myWindow.Resize(fyne.NewSize(1280, 720)) // Установим фиксированный размер окна

	employeeList = widget.NewList(
		func() int {
			return len(employees)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Имя:"),
				widget.NewLabel("   Дети:"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(employees[id].Name)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("   %d детей", len(employees[id].Children)))
		},
	)

	var addEmployeeButton *widget.Button // объявляем переменную, чтобы она была видна в main и в функции обработчика

	employeeList.OnSelected = func(id widget.ListItemID) {
		showChildren(myWindow, &employees[id], addEmployeeButton)
	}

	addEmployeeButton = widget.NewButtonWithIcon("Добавить сотрудника", theme.ContentAddIcon(), func() {
		showAddEmployeeDialog(myWindow, addEmployeeButton, func() {
			showMainScreen(myWindow, addEmployeeButton)
		})
	})

	content = container.NewBorder(nil, nil, nil, addEmployeeButton, employeeList)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func loadData() {
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading data file:", err)
		return
	}

	err = json.Unmarshal(file, &employees)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
	}
}

func saveData() {
	data, err := json.Marshal(employees)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	err = ioutil.WriteFile("data.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing data file:", err)
	}
}

func showAddEmployeeDialog(window fyne.Window, addEmployeeButton *widget.Button, callback func()) {
	nameEntry := widget.NewEntry()
	childrenEntry := widget.NewEntry()

	form := widget.NewForm(
		&widget.FormItem{Text: "Name", Widget: nameEntry},
		&widget.FormItem{Text: "Children (comma-separated)", Widget: childrenEntry},
	)

	form.OnSubmit = func() {
		childrenNames := splitAndTrim(childrenEntry.Text)
		children := make([]Child, len(childrenNames))
		for i, childName := range childrenNames {
			children[i] = Child{Name: childName}
		}

		employees = append(employees, Employee{
			Name:     nameEntry.Text,
			Children: children,
		})
		saveData()
		showMainScreen(window, addEmployeeButton)

		dialog.NewInformation("Employee Added", "Employee has been successfully added.", window).Show()
		callback()
	}

	dialog.ShowForm("Add Employee", "Submit", "Cancel", form.Items, func(bool) {}, window)
}

func showChildren(window fyne.Window, employee *Employee, addEmployeeButton *widget.Button) {
	backButton := widget.NewButtonWithIcon("Назад2", theme.NavigateBackIcon(), func() {
		showMainScreen(window, addEmployeeButton)
	})

	childrenList := widget.NewList(
		func() int {
			return len(employee.Children)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Name:"),
				widget.NewLabel("Age:"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			child := employee.Children[id]
			hBox, ok := obj.(*fyne.Container)
			if ok {
				hBox.Objects[0].(*widget.Label).SetText(child.Name)
				hBox.Objects[1].(*widget.Label).SetText(strconv.Itoa(child.Age))
			}
		},
	)

	childrenList.OnSelected = func(id widget.ListItemID) {
		showChildDetails(window, employee, addEmployeeButton, childrenList, &employee.Children[id])
	}
	childrenList.Resize(fyne.NewSize(1100, 700))
	addChildButton := widget.NewButtonWithIcon("Добавить ребенка", theme.ContentAddIcon(), func() {
		showAddChildDialog(window, employee, func() {
			loadData()
			childrenList.Refresh()
		})
	})
	addChildButton.Resize(fyne.NewSize(300, 150))
	backButton.Resize(fyne.NewSize(10, 10))

	content = container.NewBorder(nil,
		addChildButton,
		nil,
		backButton,
		childrenList,
	)

	window.SetContent(content)
}

func showChildDetails(window fyne.Window, employee *Employee, addEmployeeButton *widget.Button, childrenList *widget.List, child *Child) {
	backButton := widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {
		showChildren(window, employee, addEmployeeButton)
	})

	var content fyne.CanvasObject

	if child != nil {
		content = container.NewVBox(
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
				backButton,
				widget.NewLabel("Детали "+child.Name+":"),
			),
			widget.NewLabel("Имя: "+child.Name),
			widget.NewLabel(fmt.Sprintf("Возраст: %d", child.Age)),
			widget.NewLabel("Внешность: "+child.Appearance),
			widget.NewLabel("Любимый цвет: "+child.FavoriteColor),
			widget.NewLabel("Комментарии: "+child.Comments),
		)
	} else {
		content = widget.NewLabel("Выберите ребенка для просмотра деталей.")
	}

	window.SetContent(content)
}

func showMainScreen(window fyne.Window, addEmployeeButton *widget.Button) {
	loadData()
	employeeList.Refresh()

	content = container.NewBorder(nil, nil, nil, addEmployeeButton, employeeList)

	// Обновление содержимого главного окна
	window.SetContent(content)
}

func showAddChildDialog(window fyne.Window, employee *Employee, callback func()) {
	nameEntry := widget.NewEntry()
	ageEntry := widget.NewEntry()
	appearanceEntry := widget.NewEntry()
	colorEntry := widget.NewEntry()
	commentsEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Age", Widget: ageEntry},
			{Text: "Appearance", Widget: appearanceEntry},
			{Text: "Favorite Color", Widget: colorEntry},
			{Text: "Comments", Widget: commentsEntry},
		},
		OnSubmit: func() {
			age, err := strconv.Atoi(ageEntry.Text)
			if err != nil {
				dialog.NewError(fmt.Errorf("Invalid age"), window)
				return
			}

			employee.Children = append(employee.Children, Child{
				Name:          nameEntry.Text,
				Age:           age,
				Appearance:    appearanceEntry.Text,
				FavoriteColor: colorEntry.Text,
				Comments:      commentsEntry.Text,
			})
			saveData()

			dialog.NewInformation("Child Added", "Child has been successfully added.", window).Show()
			callback()
		},
	}

	dialog.ShowForm("Add Child", "Add", "Cancel", form.Items, func(bool) {}, window)
}

func funcIDFunc(id widget.ListItemID) fyne.CanvasObject {
	employee := employees[id]
	return container.NewHBox(
		widget.NewLabel(employee.Name),
		widget.NewLabel(fmt.Sprintf("   %d детей", len(employee.Children))),
		layout.NewSpacer(),
	)
}

func splitAndTrim(s string) []string {
	values := strings.Split(s, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}
