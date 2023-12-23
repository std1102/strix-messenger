package common

import (
	"fmt"
	"time"
)

type Pageable struct {
	PageSize   uint
	PageNumber uint
}

type fn func(*interface{}) interface{}

type ListDto struct {
	Data         []interface{}
	InputData    *[]interface{} `json:"-"`
	TotalElement uint
}

func (list *ListDto) NewListDto(input *[]interface{}, convertFn fn) *ListDto {
	list.InputData = input
	list.convert(convertFn)
	list.count()
	return list
}

func (list *ListDto) convert(convertFn fn) {
	var output []interface{}
	for i := 0; i < len(*list.InputData); i++ {
		output = append(output, convertFn)
	}
}

func (list *ListDto) count() {
	counter := 0
	for i := 0; i < len(list.Data); i++ {
		counter++
	}
	list.TotalElement = uint(counter)
}

func FormatTime(t *time.Time) string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}
