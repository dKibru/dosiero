package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

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

	// check if indexing is allowed // WHAT?
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
}

func getFileDownload(c *gin.Context) {
	var request getPathInfoRequest
	ex, err := os.Executable()
	if err != nil {
		c.JSON(404, gin.H{
			"dir":     "PAGE_NOT_FOUND",
			"message": "wrong initial directory",
		})
		return
	}
	exPath := filepath.Dir(ex)
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}
	if request.Dir == "" {
		request.Dir = exPath
	}

	newPath := filepath.Join(request.Dir, request.File)

	if !strings.HasPrefix(newPath, exPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "don't be silly"})
		return
	}

	file, err := os.Open(newPath) //Create a file
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
			Message: "File loading failed:" + err.Error(),
		})
	}
	return
}

type getPathInfoRequest struct {
	Dir  string `form:"dir" json:"dir"`
	File string `form:"file" json:"file"`
}

func getPathInfo(c *gin.Context) {
	var request getPathInfoRequest
	ex, err := os.Executable()
	if err != nil {
		c.JSON(404, gin.H{
			"dir":     "PAGE_NOT_FOUND",
			"message": "wrong initial directory",
		})
		return
	}
	exPath := filepath.Dir(ex)
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong request format"})
		return
	}
	if request.Dir == "" {
		request.Dir = exPath
	}

	newPath := filepath.Join(request.Dir, request.File)

	if !strings.HasPrefix(newPath, exPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "don't be silly"})
		return
	}

	files, err := os.ReadDir(newPath)
	if err != nil {
		c.JSON(404, gin.H{
			"dir":     exPath,
			"message": "wrong path",
		})
		return
	}
	list := []FileInfo{}
	f0 := FileInfo{
		Id:      "../",
		Name:    "..",
		Size:    0,
		ModTime: "",
		IsDir:   true,
	}
	list = append(list, f0)
	for _, entry := range files {
		fi, err := os.Stat(exPath)
		if err != nil {
			c.JSON(404, gin.H{
				"code": "PAGE_NOT_FOUND", "message": "Page not found",
			})
			return
		}
		f := FileInfo{
			Id:      entry.Name(),
			Name:    entry.Name(),
			Size:    fi.Size(),
			Mode:    fi.Mode(),
			ModTime: fi.ModTime().Format(time.RFC1123),
			IsDir:   entry.IsDir(),
		}
		list = append(list, f)
	}

	c.IndentedJSON(http.StatusOK, upcatInstance{
		Temp:           exPath,
		BaseDirectory:  newPath,
		DirectoryItems: list,
	})
}

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

	router := gin.Default()

	tbox, _ := rice.FindBox("assets")
	router.StaticFS("/assets", tbox.HTTPBox())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/api/meta", getPathInfo)
	router.POST("/api/download", getFileDownload)

	fs := EmbedFolder(server, "templates", true)
	router.Use(static.Serve("/", fs))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND", "message": "Page not found",
		})
	})
	the_ip := GetOutboundIP()
	fmt.Println("Hosted at :" + the_ip.String())
	router.Run("0.0.0.0:8080")
}
