package Tools

import (
	"fmt"
	"testing"
)

func TestGetFilesList(t *testing.T) {
	fmt.Println("----------- TestGetFilesList -----------")
	files := GetFilesList("")
	fmt.Printf("Найденные файлы: %v\n", files)
	fmt.Printf("Количество файлов: %v\n", len(files))
	//fmt.Println("----------------------------------------")
}

func TestAlbum_CreateNew(t *testing.T) {
	fmt.Println("---------- TestAlbum_CreateNew ---------")
	var a Album
	a.Description = "This doesn't seem to actually change the terminal location after I run the go program. When I print the error, I see that there was no error for os.Chdir(). I see that error is <nil>"
	a.Title = "TestAlbum"
	err := a.CreateNew()
	if err != nil {
		t.Error("Ошибка при созданиий альбома: ", err.Error())
	}
}

func TestGetAlbums(t *testing.T) {
	fmt.Println("------------ TestGetAlbums -------------")
	list := GetAlbums()
	if list == nil {
		t.Error("Внутренняя ошибка")
	}
	for _, l := range list {
		fmt.Printf("\n>>>\nНазвание: %v\nОписание: %v\nПуть: %v\nИзображения: %v\nОбложка: %v\n", l.Title, l.Description, l.Path, l.Images, l.Preview)
	}
}

func TestIsEmptyAlbum(t *testing.T) {
	fmt.Println("----------- TestIsEmptyAlbum -----------")
	res := IsEmptyAlbum("..//html/img/empty/")
	if res != true {
		t.Error("Не определяет пустой альбом: ", res)
	}
	res = IsEmptyAlbum("..//html/img/not_empty/")
	if res != false {
		t.Error("Не определяет не пустой альбом: ", res)
	}
	res = IsEmptyAlbum("..//html/img/invalid/")
	if res != true {
		t.Error("Ошибка чтения не обрабатывается: ", res)
	}
}
