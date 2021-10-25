package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// └─$ ssh nhit05@172.20.53.230                                                                                                                             130 ⨯

//go:embed assets/*
var jsassets embed.FS

//go:embed templates/*
var server embed.FS

type embedFileSystem struct {
	http.FileSystem
	indexes bool
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	f, err := e.Open(path)
	if err != nil {
		return false
	}

	// check if indexing is allowed
	s, _ := f.Stat()
	if s.IsDir() && !e.indexes {
		return false
	}

	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string, index bool) static.ServeFileSystem {
	subFS, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(subFS),
		indexes:    index,
	}
}

type upcatError struct {
	Code    int
	Message string
}
type upcatInstance struct {
	Temp           string     `json:"tmp"`
	BaseDirectory  string     `json:"base_directory"`
	DirectoryItems []FileInfo `json:"directory_items"`
}
type FileInfo struct {
	Id      string
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime string
	IsDir   bool
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "index",
	})
	// tmpl := template.Must(template.ParseFiles("layout.html"))
	// tmpl.Execute(w, data)
}

func getFileDownload(c *gin.Context) {
	customDir, _ := c.GetQuery("d")
	customPath, _ := c.GetQuery("c")
	if customDir == "" || customPath == "" {
		panic("file not found")
	}
	exPath := path.Join(customDir, customPath)

	file, err := os.Open(exPath) //Create a file
	if err != nil {
		c.JSON(http.StatusNotFound, upcatError{
			Code:    http.StatusInternalServerError,
			Message: "Nice:" + err.Error(),
		})
		return
	}
	defer file.Close()
	c.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusNotFound, upcatError{
			Code:    http.StatusInternalServerError,
			Message: "文件加载失败:" + err.Error(),
		})
	}
	return

}

func getMetaData(c *gin.Context) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	// exPath = path.Join(exPath, "/../")

	customDir, _ := c.GetQuery("d")
	customPath, _ := c.GetQuery("c")
	if customDir != "" && customPath != "" {
		exPath = path.Join(customDir, customPath)
	}
	// exPath := exPat + "/../"

	files, err := ioutil.ReadDir(exPath)
	if err != nil {
		log.Fatal(err)
	}

	list := []FileInfo{}
	f0 := FileInfo{
		Id:      "/../",
		Name:    "..",
		Size:    0,
		ModTime: "",
		IsDir:   true,
	}
	list = append(list, f0)
	for _, entry := range files {
		f := FileInfo{
			Id:      entry.Name(),
			Name:    entry.Name(),
			Size:    entry.Size(),
			Mode:    entry.Mode(),
			ModTime: entry.ModTime().Format(time.RFC1123),
			IsDir:   entry.IsDir(),
		}
		list = append(list, f)
	}

	tmp, _ := c.GetQuery("d")

	var albums = upcatInstance{
		Temp:           tmp,
		BaseDirectory:  exPath,
		DirectoryItems: list,
	}
	c.IndentedJSON(http.StatusOK, albums)
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// dsn := "root:@tcp(127.0.0.1:3306)/goeg?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// // Migrate the schema
	// db.AutoMigrate(&Product{})

	router := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	// router.LoadHTMLGlob("templates/*.html")
	// router.Static("/assets", "./assets")

	tbox, _ := rice.FindBox("assets")
	// tbox, _ := rice.MustFindBox("assets")
	router.StaticFS("/assets", tbox.HTTPBox())
	// router.StaticFS("/assets", jsassets)
	// fmt.Println(http.Dir("./assets"))
	// router.StaticFS("/assets", http.FS(jsassets))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/api/v1/meta", getMetaData)
	router.GET("/api/v1/download", getFileDownload)

	// router.GET("/", indexPage)

	fs := EmbedFolder(server, "templates", true)
	router.Use(static.Serve("/", fs))

	// fss := EmbedFolder(jsassets, "assets", false)
	// router.Use(static.Serve("/assets", fss))
	//
	// abox, _ := rice.FindBox("assets")
	// router.Use(static.Serve("/assets", abox. ))

	// fmt.Println(abox.Name())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND", "message": "Page not found",
		})
	})
	the_ip := GetOutboundIP()
	fmt.Println("Hosted at :" + the_ip.String())
	// router.Run(the_ip.String() + ":8080")
	router.Run("0.0.0.0:8080")

}
