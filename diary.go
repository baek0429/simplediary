package diary

import (
	"encoding/base64"
	"fmt"
	passwd "github.com/howeyc/gopass"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Diary struct {
	title   string
	content []byte
	date    string
}

const (
	passwordHash = "$2a$10$tZl/5h/zZMOxOqhEtMNMz.CNCMGcX/zfJwVrUPPwhtciFLzKzqKhG"
)

func generateHash(str string) []byte {
	a, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return a
}

func Run() {
	fmt.Print("keyword: ")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "encrypt":
		d := getJustWrittenDiaryFromDirectory()
		prepareDirectoryToSaveEncryptedDiary(d)
		saveDiaryToDirectory(d)
	case "decrypt":
		var date string
		var pass []byte
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(pass))
		i := 0
		for err != nil {
			i++
			fmt.Printf("Password: ")
			pass = passwd.GetPasswdMasked()
			err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(pass))
			if i > 4 {
				return
			}
		}

		fmt.Printf("Year(xxxx-xx or xxxx-xx~xxxx-xx): ")
		fmt.Scanln(&date)
		dates := yearMonthConvert(date)
		var d []*Diary
		for _, date := range *dates {
			d = *openDiaryfromDirectory(&d, date)
		}
		for _, v := range d {
			fmt.Println("decrypted", (*v).date, (*v).title)
		}
		saveDecryptedDiaryToDirectory(&d)
		return
	case "clean":
		os.RemoveAll("./decrypted")
	case "about":
		fmt.Println("Chungseok Baek, 2015.\n The program safely encrypt your diary in the relative filesystem.\nfor more instruction, type 'help'")
		Run()
	case "help":
		fmt.Println("keyword: decrypt, encrypt, about, help")
		Run()
	default:
		return
	}
}

func yearMonthConvert(date string) *[]string {
	var relativePaths []string
	if strings.Contains(date, "~") {
		splited := strings.Split(date, "~")
		runeStart := []rune(splited[0]) //xxxx-xx 48 - 0 / 57-9
		runeEnd := []rune(splited[1])
		var monthInterval rune

		monthInterval = ((runeStart[0]-48)*1000-(runeEnd[0]-48)*1000+(runeStart[1]-48)*100-(runeEnd[1]-48)*100+
			(runeStart[2]-48)*10-(runeEnd[2]-48)*10+(runeStart[3]-48)*1-(runeEnd[3]-48)*1)*12 +
			((runeStart[5]-48)*10 - (runeEnd[5]-48)*10 + (runeStart[6]-48)*1 - (runeEnd[6]-48)*1)
		monthInterval = -monthInterval

		for i := 0; i < int(monthInterval)+1; i++ {
			relativePaths = append(relativePaths, "./"+strings.Replace(string(runeStart), "-", "/", 1))
			runeStart[6] = runeStart[6] + 1
			if runeStart[5] == 49 && runeStart[6] == 51 {
				runeStart[5] = 48
				runeStart[6] = 49
				i++
				j := 0
				runeStart[3-j] = runeStart[3-j] + 1
				for runeStart[3-j] == 58 {
					runeStart[3-j] = 48
					j++
					runeStart[3-j] = runeStart[3-j] + 1
				}
			}
			if runeStart[5] == 48 && runeStart[6] == 58 {
				runeStart[5] = 49
				runeStart[6] = 49
			}
		}
	} else {

		splited := strings.Split(date, "-")
		relativePath := "./" + splited[0] + "/" + splited[1]
		relativePaths = append(relativePaths, relativePath)
	}
	return &relativePaths
}

func saveDecryptedDiaryToDirectory(diary *[]*Diary) {
	prepareDirectoryToSaveDecryptedDiary(diary)
	for _, v := range *diary {
		relativePath := "./decrypted" + v.date
		err := ioutil.WriteFile(relativePath+"/"+v.title, v.content, 0600)
		if err != nil {
			panic(err)
		}
	}
}

func prepareDirectoryToSaveDecryptedDiary(diary *[]*Diary) {
	for _, v := range *diary {
		v.date = strings.Replace(v.date, ".", "", 1)
		dirRelativePath := "./decrypted/" + v.date
		err := os.MkdirAll(dirRelativePath, 0600)
		if err != nil {
			panic(err)
		}
	}
}

func openDiaryfromDirectory(diary *[]*Diary, relativePath string) *[]*Diary {
	fInfos, err := ioutil.ReadDir(relativePath)
	if err != nil {
		fmt.Println("Couldn't find the directory: ", relativePath)
		return diary
	}
	for _, v := range fInfos {
		if v.IsDir() {
			openDiaryfromDirectory(diary, relativePath+"/"+v.Name())
		} else {
			b, err := ioutil.ReadFile(relativePath + "/" + v.Name())
			if err != nil {
				panic(err)
			}
			//./2015/12/02-title
			*diary = append(*diary, &Diary{title: v.Name() + ".txt", content: b, date: strings.Split(relativePath, "-")[0]})
			// date ./2015/12/02
			n := len(*diary)
			(*diary)[n-1] = decodeDiary((*diary)[n-1])
		}
	}
	return diary
}

func decodeDiary(diary *Diary) *Diary {
	b, err := base64.StdEncoding.DecodeString(string(diary.content))
	if err != nil {
		log.Fatal("data corruption!")
	}
	diary.content = b
	return diary
}

func saveDiaryToDirectory(diary []*Diary) {
	for _, v := range diary {
		splited := strings.Split(v.date, "-")
		dirRelativePath := splited[0] + "/" + splited[1] + "/" + strings.Split(splited[2], "T")[0]
		ioutil.WriteFile(dirRelativePath+"-"+v.title, []byte(contentEncoding(v.content)), 0600)
	}
}

func contentEncoding(b []byte) string {
	str := base64.StdEncoding.EncodeToString(b)
	return str
}

func prepareDirectoryToSaveEncryptedDiary(diary []*Diary) {
	for _, v := range diary {
		splited := strings.Split(v.date, "-")
		dirRelativePath := splited[0] + "/" + splited[1]
		err := os.MkdirAll(dirRelativePath, 0600)
		if err != nil {
		}
	}
}

func getJustWrittenDiaryFromDirectory() []*Diary {
	var diary []*Diary
	f, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}
	for _, v := range f {
		if !v.IsDir() {
			nameWithExt := strings.Split(v.Name(), ".")
			if len(nameWithExt) > 1 {
				if nameWithExt[1] == "txt" { // if filename is *.txt
					byteContent, err := ioutil.ReadFile(v.Name())
					if err != nil {
						panic(err)
					}
					d := &Diary{
						date: v.ModTime().Format("2006-01-02T15-04"),
					}
					d.title = nameWithExt[0]
					d.content = byteContent
					diary = append(diary, d)
					os.RemoveAll(v.Name())
				}
			}
		}
	}
	return diary
}
