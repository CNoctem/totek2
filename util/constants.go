package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"totek2/keys"
)

type User struct {
	Username string
	Password string
}

func Authenticate() User {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username:")
	usr, _ := reader.ReadString('\n')
	usr = strings.Replace(usr, "\n", "", -1)
	pwd := keys.Askpass("Password:")
	return User{usr, pwd}
}

func AskForConfirmation(msg string) bool {
	var response string

	fmt.Print(msg + " y/n:")

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Yes or no?:")
		return AskForConfirmation(msg)
	}
}

