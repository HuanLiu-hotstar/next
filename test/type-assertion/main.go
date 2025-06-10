package main

import (
	"fmt"
)

type Data struct {
	X int
}

func getData(x int) interface{} {
	return &WidgetBinderOutput{
		TypeInstanceName: "hello",
		Data:             &Data{x},
	}
}

type WidgetBinderOutput struct {
	TypeInstanceName string
	Data             interface{}
}

func (w *WidgetBinderOutput) ToBasic() interface{} {
	if w == nil {
		return nil
	}
	return map[string]interface{}{
		"TypeInstanceName": w.TypeInstanceName,
		"Data":             w.Data,
	}
}
func (w *WidgetBinderOutput) FromBasic(input interface{}) *WidgetBinderOutput {
	inputMap := input.(map[string]interface{})
	w.TypeInstanceName, _ = inputMap["TypeInstanceName"].(string)
	w.Data = inputMap["Data"]
	return w
}

func main() {
	d := getData(10)
	a, ok := d.(*WidgetBinderOutput)
	fmt.Printf("%T ok:%t a:%+v", d, ok, a)
}
