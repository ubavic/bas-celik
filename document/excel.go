package document

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/ubavic/bas-celik/localization"
	"github.com/xuri/excelize/v2"
)

func CreateExcel(document any) ([]byte, error) {
	structType := reflect.TypeOf(document)
	structVal := reflect.ValueOf(document)

	if structType.Kind() != reflect.Struct {
		return nil, errors.New("not a struct")
	}

	f := excelize.NewFile()

	err := f.SetColWidth("Sheet1", "A", "A", 30)
	if err != nil {
		return nil, err
	}

	err = f.SetColWidth("Sheet1", "B", "B", 60)
	if err != nil {
		return nil, err
	}

	currentRow := 1
	putData := func(label, value string) {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow), label)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow), value)
		currentRow += 1
	}

	fields := reflect.VisibleFields(structType)

	for _, field := range fields {
		switch field.Type.Kind() {
		case reflect.String:
			putData(field.Name, structVal.FieldByName(field.Name).String())
		case reflect.Bool:
			str := localization.FormatYesNo(structVal.FieldByName(field.Name).Bool(), localization.En)
			putData(field.Name, str)
		}
	}

	buffer := bytes.Buffer{}

	err = f.Write(&buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
