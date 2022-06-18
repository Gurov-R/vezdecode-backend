package cli

import (
	vezdecodebackend "Gurov-R/vezdecode-backend"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	asciiMeme = `
⠀⣞⢽⢪⢣⢣⢣⢫⡺⡵⣝⡮⣗⢷⢽⢽⢽⣮⡷⡽⣜⣜⢮⢺⣜⢷⢽⢝⡽⣝
⠸⡸⠜⠕⠕⠁⢁⢇⢏⢽⢺⣪⡳⡝⣎⣏⢯⢞⡿⣟⣷⣳⢯⡷⣽⢽⢯⣳⣫⠇
⠀⠀⢀⢀⢄⢬⢪⡪⡎⣆⡈⠚⠜⠕⠇⠗⠝⢕⢯⢫⣞⣯⣿⣻⡽⣏⢗⣗⠏⠀
⠀⠪⡪⡪⣪⢪⢺⢸⢢⢓⢆⢤⢀⠀⠀⠀⠀⠈⢊⢞⡾⣿⡯⣏⢮⠷⠁⠀⠀
⠀⠀⠀⠈⠊⠆⡃⠕⢕⢇⢇⢇⢇⢇⢏⢎⢎⢆⢄⠀⢑⣽⣿⢝⠲⠉⠀⠀⠀⠀
⠀⠀⠀⠀⠀⡿⠂⠠⠀⡇⢇⠕⢈⣀⠀⠁⠡⠣⡣⡫⣂⣿⠯⢪⠰⠂⠀⠀⠀⠀
⠀⠀⠀⠀⡦⡙⡂⢀⢤⢣⠣⡈⣾⡃⠠⠄⠀⡄⢱⣌⣶⢏⢊⠂⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢝⡲⣜⡮⡏⢎⢌⢂⠙⠢⠐⢀⢘⢵⣽⣿⡿⠁⠁⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠨⣺⡺⡕⡕⡱⡑⡆⡕⡅⡕⡜⡼⢽⡻⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣼⣳⣫⣾⣵⣗⡵⡱⡡⢣⢑⢕⢜⢕⡝⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⣴⣿⣾⣿⣿⣿⡿⡽⡑⢌⠪⡢⡣⣣⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⡟⡾⣿⢿⢿⢵⣽⣾⣼⣘⢸⢸⣞⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠁⠇⠡⠩⡫⢿⣝⡻⡮⣒⢽⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

    NO MEMES?

`
	helpText = `
Список доступных команд:
  help                          - выводит все доступные команды
  exit                          - выйти из терминала (учтите, что сервер перестанет работать)
  load_vezdekod                 - загрузить мемы из группы "Вездекод"
  load_group [ссылка на группу] - загрузить мемы из другой группы ( например https://vk.com/abstract_memes )
  feed                          - лента с мемами

`
)

func RunCli() {
	fmt.Println("\n\nДобро пожаловать в терминал с мемами!\nВведите \"help\" для просмотра всех команд.")
	fmt.Print(asciiMeme)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command := readLine(*reader)
		if command != "" {
			runCommand(command, *reader)
		}
	}
}

func getLocalDomain() string {
	return "http://" + viper.GetString("host") + ":" + viper.GetString("port")
}

func readLine(reader bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.Replace(line, "\n", "", -1))
}

func runCommand(command string, reader bufio.Reader) {
	if command == "help" {
		help()
	} else if command == "exit" {
		exit()
	} else if command == "load_vezdekod" {
		load_vezdekod()
	} else if strings.HasPrefix(command, "load_group") {
		if len(strings.Split(command, " ")) < 2 {
			fmt.Println("Необходимо ввести ссылку на группу\nНапример: load_groups https://vk.com/abstract_memes")
		} else {
			load_group(strings.Split(command, " ")[1], reader)
		}
	} else if command == "feed" {
		feed(0, reader)
	} else {
		fmt.Println("Вы ввели неверную комманду. Напиши help для просмотра всех команд")
	}
}

func help() {
	fmt.Print(helpText)
}

func exit() {
	fmt.Println("Пока(")
	os.Exit(0)
}

func load_vezdekod() {
	fmt.Println("Загрузка...")
	_, err := http.PostForm(getLocalDomain()+"/api/memes/load-vezdekod", url.Values{})

	if err != nil {
		logrus.Fatalf("error on request /api/memes/load-vezdekod [POST]: %s", err)
	}

	fmt.Println("Вы успешно загрузили мемы с вездекода")
}

func load_group(address string, reader bufio.Reader) {

	fmt.Print("Введите пароль: ")
	password := readLine(reader)

	fmt.Println()

	values := map[string]string{"address": address, "password": password}
	json_data, _ := json.Marshal(values)

	fmt.Println("Загрузка...")

	resp, err := http.Post(getLocalDomain()+"/api/memes/load-group", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		logrus.Fatalf("error on request /api/memes/load-group [POST]: %s", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Неверный пароль (почитайте README.md)")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Вы не ввели пароль")
	} else {
		fmt.Println("Вы успешно загрузили мемы " + fmt.Sprint(address))
	}
}

func feed(page int, reader bufio.Reader) {

	values := map[string]int{"page": page}
	json_data, _ := json.Marshal(values)

	resp, err := http.Post(getLocalDomain()+"/api/memes/feed", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		logrus.Fatalf("error on request /api/memes/feed [POST]: %s", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var meme vezdecodebackend.Meme

	err = json.Unmarshal(body, &meme)

	if err != nil {
		logrus.Fatalf("error while unmarshling meme body: %s", err)
	}

	fmt.Printf("Вот вам мем: %s\nКоличество лайков: %d\nВыложен: %s\n", meme.ImageUrl, meme.LikesCount, time.Unix(int64(meme.Timestamp), 0))

	fmt.Printf("\n\nВведите цифру для последующих действий:\n1 - Лайк\n2 - Продвигать этот мем (будет показываться чаще остальных)\n3 - Следующий мем\n(любая другая команда) - Выйти в консоль с мемами\n")

	for {
		command := readLine(reader)

		if command == "1" {
			likeMeme(meme.Id)
		} else if command == "2" {
			promoteMeme(meme.Id, reader)
		} else if command == "3" {
			feed(page+1, reader)
			break
		} else {
			break
		}
	}

}

func promoteMeme(id int, reader bufio.Reader) {

	fmt.Print("Введите пароль: ")

	password := readLine(reader)
	fmt.Println()

	values := map[string]string{"meme_id": fmt.Sprint(id), "password": password}
	json_data, _ := json.Marshal(values)

	resp, err := http.Post(getLocalDomain()+"/api/memes/promote", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		logrus.Fatalf("error sending /promote: %s", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Неверный пароль")
	} else {
		fmt.Println("Вы продвигаете пост!")
	}
}

func likeMeme(id int) {

	values := map[string]string{"meme_id": fmt.Sprint(id)}
	json_data, _ := json.Marshal(values)

	resp, err := http.Post(getLocalDomain()+"/api/memes/like", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		logrus.Fatalf("error sending /like: %s", err)
	}

	if resp.StatusCode == http.StatusConflict {
		fmt.Println("Вы уже лайкнули этот мем")
	} else {
		fmt.Println("Вы поставили лайк!")
	}
}
