package Tools

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

const (
	DataBaseHost  = "localhost:3306"
	DataBaseName  = "akaiphoto"
	AdminPassword = "1qaz@WSX3edc$RFV"
)

type DBAlbum struct {
	Name        string
	Description string
	AccessLvl   int
	Images      []DBImage
}

func (a DBAlbum) String() string {
	return fmt.Sprintf("%s[Album: %s (%d)]%s[%d] %s[%s]", ColorGreen, a.Name, len(a.Images), ColorRed, a.AccessLvl, ColorYellow, a.Description)
}

type DBImage struct {
	Name        string
	Description string
	AccessLvl   int
}

type DBUser struct {
	Name       string `json:"name"`
	PublicName string `json:"public_name"`
	Password   string `json:"password"`
	Access     int    `json:"access_lvl"`
	Active     bool   `json:"-"`
}

func DBGetUser(name string) (DBUser, error) {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	request := fmt.Sprintf("select * from users where login = '%v'", name)
	rows, err := db.Query(request)
	if err != nil {
		fmt.Println(err.Error())
		return DBUser{}, err
	}
	defer rows.Close()
	post := &DBUser{}
	rows.Next()
	err = rows.Scan(&post.Name, &post.PublicName, &post.Password, &post.Access, &post.Active)
	if err != nil {
		fmt.Println(err.Error())
	}
	return *post, err
}

func DBUpdateUser(user DBUser) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[UpdateUser] panic occurred:", err)
		}
	}()
	request := "update users set "
	if user.PublicName != "" {
		request += fmt.Sprintf("public_name='%v', ", user.PublicName)
	}
	if user.Password != "" {
		request += fmt.Sprintf("password='%v', ", user.Password)
	}
	request += fmt.Sprintf("active = %v  where login = '%v'", user.Active, user.Name)
	_, err := db.Query(request)
	if err != nil {
		fmt.Println("[SQL] Update user: ", err.Error())
		return err
	}

	return nil
}

func DBNewUser(name, publicName, password string, access int) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[NewUser] panic occurred:", err)
		}
	}()
	if publicName == "" {
		publicName = name
	}
	request := fmt.Sprintf("insert into users values ('%v', '%v', '%v', %v, 0);", name, publicName, password, access)
	_, err := db.Query(request)
	if err != nil {
		fmt.Println("[SQL] Add user: ", err.Error())
		return err
	}
	return nil
}

func DBRemoveUser(name string) {
	request := fmt.Sprintf("delete from users where login='%v';", name)
	_, err := db.Query(request)
	if err != nil {
		fmt.Println("[SQL] Add user: ", err.Error())
		return
	}
	fmt.Println("[SQL] DBUser removed: ", name)
}

func DBCreate() {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic occurred", err)
		}
	}()
	dsn := "root@tcp(" + DataBaseHost + ")/?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		Log("SQL Open DataBase", err)
		return
	}
	//db.SetConnMaxIdleTime(10)
	request := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS " + DataBaseName + " CHARACTER SET utf8 COLLATE utf8_general_ci;")
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create database:", err)
	} else {
		NamedMessage("SQL", "Create database: Ok")
	}
}

func DBOpen() error {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic occurred", err)
		}
	}()
	dsn := "root@tcp(" + DataBaseHost + ")/" + DataBaseName + "?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		Log("SQL Open DataBase", err)
		return err
	}
	//db.SetConnMaxIdleTime(10)
	err = db.Ping()
	if err != nil {
		Log("SQL Ping DataBase", err)
		return err
	}
	return nil
}

func DBInit() {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic occurred", err)
		}
	}()
	/*
		request := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS " + DataBaseName + ";")
		_, err := db.Query(request)
		if err != nil {
			NamedMessage("SQL", "Create database:", err)
		} else {
			NamedMessage("SQL", "Create database: Ok")
		}
	*/
	request := fmt.Sprintf("create table if not exists users(login char(15) not null, public_name char(15), " +
		"password char(25), accesslvl int not null, active bool not null, primary key(login, public_name));")
	_, err := db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create table 'users':", err)
	} else {
		NamedMessage("SQL", "Create table 'users': Ok")
	}

	request = fmt.Sprintf("insert into users values('admin', 'Администратор', '" + AdminPassword + "', 0, 0);")
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create user 'admin':", err)
	} else {
		NamedMessage("SQL", "Create user 'admin': Ok")
	}

	request = fmt.Sprintf("insert into users values('guest', 'Гость', '', 10, 0);")
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create user 'guest':", err)
		//return err
	} else {
		NamedMessage("SQL", "Create user 'guest': Ok")
	}

	request = fmt.Sprintf("create table if not exists albums(name nvarchar(128) not null, description nvarchar(256), accesslvl int not null, primary key(name));")
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create table 'albums':", err)
	} else {
		NamedMessage("SQL", "Create table 'albums': Ok")
	}

	request = fmt.Sprintf("create table if not exists images(name nvarchar(128) not null, description nvarchar(256), accesslvl int not null, primary key(name));")
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create table 'images':", err)
	} else {
		NamedMessage("SQL", "Create table 'images': Ok")
	}

	request = fmt.Sprintf("create table if not exists album_image(album nvarchar(128) not null, image nvarchar(128) not null, " +
		"foreign key(album) references albums(name), foreign key(image) references images(name));")
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create table 'album_image':", err)
	} else {
		NamedMessage("SQL", "Create table 'album_image': Ok")
	}
}

func DBCreateAlbum(name string, description string, accesslvl int) {
	request := fmt.Sprintf("insert into albums values(N'%v', N'%v', %v);", name, description, accesslvl)
	_, err := db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Create album:", err)
		//return err
	} else {
		NamedMessage("SQL", "Create album: Ok")
	}
}

func DBDeleteAlbum(album DBAlbum) error {
	request := fmt.Sprintf("select * from albums where name='%v' and accesslvl >= %v", album.Name, album.AccessLvl)
	rows, err := db.Query(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer rows.Close()
	post := DBAlbum{}
	if rows.Next() {
		err = rows.Scan(&post.Name, &post.Description, &post.AccessLvl)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		return errors.New("[SQL] Альбом не найден либо нет доступа")
	}

	request = fmt.Sprintf("delete from images where name='%v';", album.Name)
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Не удалось удалить альбом", err)
		return err
	}
	return nil
}

func DBGetAlbums(accesslvl int) []DBAlbum {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	request := fmt.Sprintf("select * from albums where accesslvl >= %v", accesslvl)
	rows, err := db.Query(request)
	if err != nil {
		fmt.Println(err.Error())
		NamedMessage("SQL", "", err)
		return nil
	}
	defer rows.Close()
	post := DBAlbum{}
	res := make([]DBAlbum, 0, 10)
	for rows.Next() {
		err = rows.Scan(&post.Name, &post.Description, &post.AccessLvl, &post.Images)
		if err != nil {
			NamedMessage("SQL", "", err)
			continue
		}
		res = append(res, post)
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func DBGetAlbum(album DBAlbum) (DBAlbum, error) {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	request := fmt.Sprintf("select * from albums where name='%s' and accesslvl >= %d", album.Name, album.AccessLvl)
	rows, err := db.Query(request)
	if err != nil {
		return album, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&album.Name, &album.Description, &album.AccessLvl, &album.Images)
		if err != nil {
			return album, err
		}
	}

	return album, nil
}

func DBAddImage(image DBImage) error {
	request := fmt.Sprintf("insert into images values('%v', '%v', %v);", image.Name, image.Description, image.AccessLvl)
	_, err := db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Не удалось добавить изображение ", image.Name, " в базу: ", err)
		return err
	}
	return nil
}

func DBGetImage(image DBImage) (DBImage, error) {
	defer func() {
		if err := recover(); err != nil {
			Log("Panic", err)
		}
	}()
	request := fmt.Sprintf("select * from images where name='%s' and accesslvl >= %d", image.Name, image.AccessLvl)
	rows, err := db.Query(request)
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
		return DBImage{}, err
	}
	post := DBImage{}
	if rows.Next() {
		err = rows.Scan(&post.Name, &post.Description, &post.AccessLvl)
		if err != nil {
			return DBImage{}, err
		}
	}
	return post, nil
}

func DBAddImageToAlbum(image DBImage, album DBAlbum) error {
	_, err := DBGetImage(image)
	if err != nil {
		return err
	}
	request := fmt.Sprintf("insert into album_image values('%v', '%v');", album.Name, image.Name)
	_, err = db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Не удалось добавить изображение ", image.Name, " в альбом: ", err)
		return err
	}

	return nil
}

func DBDeleteImage(image DBImage) error {
	request := fmt.Sprintf("delete from images where name='%v';", image.Name)
	_, err := db.Query(request)
	if err != nil {
		NamedMessage("SQL", "Не удалось удалить изображение ", image.Name, err)
		return err
	}

	DeleteImage(image.Name)
	return nil
}
