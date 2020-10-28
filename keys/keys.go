package keys

import (
	"fmt"
	"github.com/eiannone/keyboard"
)

func Askpass(msg string) string {
	pwd := make([]rune, 0)
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Print(msg)
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		pwd = append(pwd, event.Rune)
		if event.Key == keyboard.KeyEnter {
			break
		}
	}
	fmt.Println()
	return string(pwd)
}
