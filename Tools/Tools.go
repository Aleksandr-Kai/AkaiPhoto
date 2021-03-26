package Tools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	RootDir     = "./html/img/"
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type TError struct {
	Code    int16
	Message string
}

type Album struct {
	Path        string   `json:"-"`
	Title       string   `json:"title""`
	Description string   `json:"description""`
	Preview     string   `json:"image"`
	Images      []string `json:"-"`
	Thumbs      []string `json:"-"`
}

func (p *Album) CreateNew() error {
	aName := time.Now().Unix()
	p.Path = strconv.FormatInt(aName, 10)
	os.Mkdir(RootDir+p.Path, os.ModePerm)
	desc, err := os.Create(RootDir + p.Path + "/description.txt")
	if err != nil {
		fmt.Printf("Не удалось создать файл описания альбома: %v\n", err.Error())
		return err
	}
	enc, _ := json.Marshal(p)
	_, err = desc.Write(enc)
	desc.Close()
	return err
}

func (p *Album) FullPath() string {
	return "img/" + p.Path + "/"
}

func IsEmptyAlbum(path string) bool {
	files, err := ioutil.ReadDir(path) // читаем папку с альбомами
	if err != nil {
		fmt.Printf("[IsEmptyAlbum]: %v\n", err.Error())
		return true
	}
	for _, file := range files {
		if !file.IsDir() && (strings.ToLower(filepath.Ext(file.Name())) == ".jpg") {
			return false
		}
	}
	return true
}

func GetAlbums() []Album {
	files, err := ioutil.ReadDir(RootDir) // читаем папку с альбомами
	if err != nil {
		fmt.Errorf("[GetAlbums]: %v", err.Error())
		return nil
	}
	albums := make([]Album, 0, 10)
	for _, file := range files {
		if file.IsDir() { // если папка, т.е. альбом
			txt, err := ioutil.ReadFile(RootDir + file.Name() + "/description.txt") // пытаемся читать описание
			al := &Album{}
			if err != nil {
				fmt.Printf("[GetAlbums] Не удалось открыть файл описания: %v\n", err.Error())
			}
			eerr := json.Unmarshal(txt, al)
			if eerr != nil {
				fmt.Printf("[GetAlbums] Ошибка распаковки json: %v\n", eerr.Error())
			}
			al.Path = file.Name()
			if (err != nil) || (eerr != nil) { // если не удалось, пытаемся исправить
				al.Title = "Album_Title"
				al.Description = "Album_Description"
				dFile, err := os.Create(RootDir + al.Path + "/description.txt")
				if err != nil { // НЕ ВАЛИДНЫЙ АЛЬБОМ
					fmt.Printf("[GetAlbums] Не удалось создать файл описания альбома: %v\n", err.Error())
					continue
				} else {
					txt, _ = json.Marshal(al)
					_, err = dFile.Write(txt)
					dFile.Close()
					if err != nil { // НЕ ВАЛИДНЫЙ АЛЬБОМ
						continue
					}
				}
			}

			images, err := ioutil.ReadDir(RootDir + file.Name() + "/s/") // читаем папку с альбомами
			if err == nil {
				for _, img := range images {
					if !img.IsDir() && (strings.ToLower(filepath.Ext(img.Name())) == ".jpg") {
						al.Images = append(al.Images, img.Name())
					}
				}
			}
			if (al.Preview == "") && (len(al.Images) > 0) {
				al.Preview = al.Images[0]
			}

			al.Description += " \n" + file.Name()
			albums = append(albums, *al)
		}
	}
	return albums
}

func GetFilesList(path string) []string {
	files, err := ioutil.ReadDir(RootDir + path)
	if err != nil {
		fmt.Printf("Read File Sistem Fail: %v\n", err.Error())
		return nil
	}
	res := make([]string, 0)
	var dir []string
	for _, file := range files {
		if file.IsDir() {
			dir = GetFilesList(path + file.Name() + "/")
			if dir == nil {
				continue
			}
			res = append(res, dir...)
		} else {
			if strings.ToLower(filepath.Ext(file.Name())) == ".jpg" {
				res = append(res, path+file.Name())
			}
		}
	}
	return res
}

func GetTmpImages() (error, []string) {
	files, err := ioutil.ReadDir(RootDir + "upload/")
	if err != nil {
		Log("Read File Sistem Fail", err)
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		fmt.Println(ex)
		return err, nil
	}
	res := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			if strings.ToLower(filepath.Ext(file.Name())) == ".jpg" {
				res = append(res, file.Name())
			}
		}
	}
	return nil, res
}

func Log(msg string, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[Log] Epic fail!")
		}
	}()
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	tmp := strings.Split(funcName, "/")
	funcName = tmp[len(tmp)-1]
	str := fmt.Sprintf("[%v] %v", funcName, msg)
	if len(args) > 0 {
		str += fmt.Sprintf(":")
	}
	for _, arg := range args {
		str += fmt.Sprintf(" %v", arg)
	}

	log.Println(str)
}

func Message(msg string, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	str := fmt.Sprintf("[%v] %v", funcName, msg)
	if len(args) > 0 {
		str += fmt.Sprintf(":")
	}
	for _, arg := range args {
		str += fmt.Sprintf(" %v", arg)
	}

	fmt.Println(str)
}

func NamedMessage(prefix string, args ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	if len(args) == 0 {
		Log("No arguments")
		return
	}
	str := fmt.Sprintf("[%v]", prefix)
	for _, arg := range args {
		str += fmt.Sprintf(" %v", arg)
	}

	fmt.Println(str)
}

func GetHash(file *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)[:16]), nil
}

// Processing uploaded images
// Make S, M and L sizes. Place it to S, M and L folders
func ImgsProcess(album string) {
	// Получение списка файлов для обработки
	err, imgs := GetTmpImages()
	if err != nil {
		Log(err.Error())
		return
	}
	var hashName string // Уникальне имя для изображения
	for _, f := range imgs {
		path := RootDir + "upload/" + f
		Log("Обработка", path)
		imgIn, err := os.Open(path)
		if err != nil {
			Log("Не удалось открыть изображение", err)
			continue
		}
		if hashName, err = GetHash(imgIn); err == nil {
			Log("Новое имя файлов", hashName)
			imgIn.Seek(0, 0)
			imgJpg, err := jpeg.Decode(imgIn)
			if err != nil {
				Log("Ошибка декодирования", err)
			} else {
				imgS := resize.Thumbnail(400, 400, imgJpg, resize.Bicubic)
				imgM := resize.Thumbnail(1920, 1080, imgJpg, resize.Bicubic)

				if _, err := os.Stat(RootDir + "s/" + hashName + ".jpg"); !os.IsNotExist(err) {
					imgIn.Close()
					Log("Изображение с таким именем уже существует", hashName)
				} else {
					imgOut, err := os.Create(RootDir + "s/" + hashName + ".jpg")
					if err != nil {
						Log("Не удалось сохранить S файл", err)
					} else {
						jpeg.Encode(imgOut, imgS, nil)
						imgOut.Close()
					}
				}

				if _, err := os.Stat(RootDir + "m/" + hashName + ".jpg"); !os.IsNotExist(err) {
					imgIn.Close()
					Log("Изображение с таким именем уже существует", hashName)
				} else {
					imgOut, err := os.Create(RootDir + "m/" + hashName + ".jpg")
					if err != nil {
						Log("Не удалось сохранить M файл", err)
					} else {
						jpeg.Encode(imgOut, imgM, nil)
						imgOut.Close()
					}
				}
			}
		}
		imgIn.Close()
		if _, err := os.Stat(RootDir + "l/" + hashName + ".jpg"); !os.IsNotExist(err) {
			imgIn.Close()
			Log("Изображение с таким именем уже существует", hashName)
		} else {
			err = os.Rename(path, RootDir+"l/"+hashName+".jpg")
			if err != nil {
				Log("Не удалось сохранить L файл", err)
			}
		}
		DBAddImage(DBImage{hashName, "", 10})
		DBAddImageToAlbum(DBImage{Name: hashName, AccessLvl: 10}, DBAlbum{Name: album, AccessLvl: 10})
		Log("*********************************")
	}
}

func DeleteImage(name string) {
	err := os.Remove(RootDir + "upload/s/" + name + ".jpg")
	if err != nil {
		Log(err.Error())
	}
	err = os.Remove(RootDir + "upload/m/" + name + ".jpg")
	if err != nil {
		Log(err.Error())
	}
	err = os.Remove(RootDir + "upload/l/" + name + ".jpg")
	if err != nil {
		Log(err.Error())
	}
	err = os.Remove(RootDir + "upload/upload/" + name + ".jpg")
	if err != nil {
		Log(err.Error())
	}
}
