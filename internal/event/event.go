package event

import "fmt"

func EventListener(ch chan interface{}) {
	for {
		val, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(val)
	}
}
