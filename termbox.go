package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			default:
				fmt.Printf("Char string: %s\n", string(ev.Ch))
			}

			//printf_tb(8, 22, termbox.ColorRed, termbox.ColorBlack, "string:  %s", funckeymap(ev.Key))
			//printf_tb(60, 22, termbox.ColorRed, termbox.ColorBlack, "string:  %s", string(ev.Ch))
		case termbox.EventError:
			panic(ev.Err)
		}
		//termbox.Flush()
	}
}
