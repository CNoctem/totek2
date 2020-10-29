package util

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"totek2/keys"
)

type User struct {
	Username string
	Password string
}

func Authenticate() User {
	log.Println("Authenticating...")
	saveduser, err := ReadUserData()
	if err == nil {
		log.Println("Found user data in config.")
		return *saveduser
	} else {
		log.Println("userdata not found")
		log.Println(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username:")
	usr, _ := reader.ReadString('\n')
	usr = strings.Replace(usr, "\n", "", -1)
	pwd := keys.Askpass("Password:")
	pwd = string(bytes.Trim([]byte(pwd), "\x00"))
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

func URLWithUser(url string, user User) (*string, error) {
	idx := strings.Index(url, "://")
	if idx == -1 {
		return nil, errors.New("malformed url")
	}
	newUrl := url[:idx + 3] + user.Username + ":" + user.Password + "@" + url[idx+3:]
	return &newUrl, nil
}

func SaveUserData(tUser User) {
	home, err := getConfigFile()
	if _, err := os.Stat(home); os.IsNotExist(err) {
		err = os.Mkdir(home, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.Create(filepath.Join(home, "userdata"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(tUser.Username + "\n" + tUser.Password)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Println(err)
	}
}

func ReadFile(path string) (*[]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		fmt.Printf(" > Failed with error %v\n", scanner.Err())
		return nil, scanner.Err()
	}
	return &lines, nil
}

func ReadUserData() (*User, error) {
	cf, err := getConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	lines, err := ReadFile(filepath.Join(cf, "userdata"))
	if err != nil {
		return nil, err
	}
	return &User{Username: (*lines)[0], Password: (*lines)[1]}, nil
}

func DirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func TExit(code int, msg... string) {
	for _, m := range msg {
		log.Println(m)
	}
	if code < 0 || code > len(msgs) - 1 {
		code = 0
	}
	fmt.Println("\n\n---------------------------------------------------------\n\t", msgs[code])
	fmt.Println()
	os.Exit(code)
}

func TExitOnError(err error, code... int) {
	if err != nil {
		c := 4
		if len(code) > 0 {
			c = code[0]
		}
		log.Println(err)
		TExit(c)
	}
}

func CreateDirOrExit(path string) {
	err := os.Mkdir(path, 0755)
	TExitOnError(err)
}

func getConfigFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	home := filepath.Join(usr.HomeDir, ".totek2")
	return home, err
}

var msgs = []string{
	"- Háromba vágtad, édes, jó Lajosom?",
	"Ha még egyszer be meri tenni a lábát, legyen olyan kedves Mariskám, lője agyon.",
	"Hát ez csak természetes, mélyen tisztelt Őrnagy úr!",
	"Kilenc hónap állandó zajban és büdösségben eltöltött frontszolgálat után...",
	"Mit csinál a kedves Papa abban a büdös budiban???"}


