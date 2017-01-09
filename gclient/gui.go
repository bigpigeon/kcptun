package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/andlabs/ui"
)

// full-bit integer
const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func fillSelectField() {

}

func makerule(s string) map[string]string {
	rulemap := map[string]string{}
	for _, v := range strings.Split(s, ";") {
		if v != "" {
			var key, value string
			for i, korv := range strings.SplitN(v, ":", 2) {
				if i == 0 {
					key = korv
				} else {
					value = korv
				}
			}
			rulemap[key] = value
		}
	}
	return rulemap
}

func fillConfigPage(c *Config, tab *ui.Tab, reloadCallback func(*ui.Button)) {
	value := reflect.ValueOf(c).Elem()

	boxgroup := [2]*ui.Box{ui.NewVerticalBox(), ui.NewVerticalBox()}

	controls := []func() interface{}{}
	for i := 0; i < value.NumField(); i++ {
		box := boxgroup[i%2]
		field := value.Type().Field(i)

		label := field.Name

		dvalue := value.Field(i).Interface()
		box.Append(ui.NewLabel(label), false)

		rulemap := makerule(field.Tag.Get("ui"))
		switch field.Type.Kind() {
		case reflect.String:
			if rulemap["select"] != "" {
				options := strings.Split(rulemap["select"], ",")
				combo := ui.NewCombobox()
				selected := 0
				for index, op := range options {
					if op == dvalue.(string) {
						selected = index
					}
					combo.Append(op)
				}
				combo.SetSelected(selected)
				box.Append(combo, false)
				controls = append(controls, func() interface{} {
					return options[combo.Selected()]
				})
			} else {
				entry := ui.NewEntry()
				entry.SetText(dvalue.(string))
				box.Append(entry, false)
				controls = append(controls, func() interface{} {
					return entry.Text()
				})
			}
		case reflect.Int:

			if rulemap["select"] != "" {
				var options []int
				for _, s := range strings.Split(rulemap["select"], ",") {
					d, err := strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
					options = append(options, d)
				}
				combo := ui.NewCombobox()
				selected := 0
				for index, op := range options {
					if op == dvalue.(int) {
						selected = index
					}
					combo.Append(strconv.Itoa(op))
				}
				combo.SetSelected(selected)
				box.Append(combo, false)
				controls = append(controls, func() interface{} {
					return options[combo.Selected()]
				})
			} else {
				var min, max int
				fmt.Sscanf(rulemap["range"], "%d-%d", &min, &max)
				if min == 0 && max == 0 {
					min, max = MinInt, MaxInt
				}
				spinbox := ui.NewSpinbox(min, max)
				spinbox.SetValue(dvalue.(int))
				box.Append(spinbox, false)
				controls = append(controls, func() interface{} {
					return spinbox.Value()
				})
			}
		case reflect.Bool:
			boolselect := ui.NewCombobox()
			boolselect.Append("false")
			boolselect.Append("true")
			if dvalue.(bool) == true {
				boolselect.SetSelected(1)
			}
			box.Append(boolselect, false)
			controls = append(controls, func() interface{} {
				return boolselect.Selected() == 1
			})
		default:
			panic("ui datatype not found" + field.Type.Name())
		}
	}
	parent := ui.NewHorizontalBox()
	parent.Append(boxgroup[0], false)
	parent.Append(boxgroup[1], false)

	enter := ui.NewButton("确认")
	enter.OnClicked(func(b *ui.Button) {
		for i := 0; i < value.NumField(); i++ {

			value.Field(i).Set(reflect.ValueOf(controls[i]()))
		}
		err := saveJsonConfig(c, ConfigFile)
		if err != nil {
			panic(err)
		}
		reloadCallback(b)
	})

	pparent := ui.NewVerticalBox()
	pparent.Append(parent, false)
	pparent.Append(enter, false)
	tab.Append("config", pparent)
}

func CreateUi(c *Config) {
	err := ui.Main(func() {
		logPage := ui.NewVerticalBox()
		tab := ui.NewTab()
		var wg sync.WaitGroup
		stop := make(chan bool)
		fillConfigPage(c, tab, func(b *ui.Button) {
			stop <- true
			wg.Wait()
			go background(c, stop, &wg)
		})
		tab.Append("log", logPage)

		window := ui.NewWindow("kcptun", 200, 300, false)
		window.SetChild(tab)
		window.Show()
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		go background(c, stop, &wg)
	})
	if err != nil {
		panic(err)
	}
}
