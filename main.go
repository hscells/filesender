package main

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

type config struct {
	FileName string `toml:"fileName"`
	Port     string `toml:"port"`
	Host     string `toml:"host"`
}

var (
	keys map[string]*bufio.Scanner
)

func main() {
	var config config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	keys = make(map[string]*bufio.Scanner)

	r := gin.Default()

	r.GET("/new", func(c *gin.Context) {
		id := uuid.New().String()
		f, err := os.OpenFile(config.FileName, os.O_RDONLY, 0664)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)
		keys[id] = scanner
		c.JSON(http.StatusOK, gin.H{
			"key": id,
		})
		return
	})

	r.GET("/line/:key", func(c *gin.Context) {
		key := c.Param("key")
		if _, ok := keys[key]; !ok {
			log.Println("key not found")
			c.Status(http.StatusInternalServerError)
			return
		}
		scanner := keys[key]
		if scanner.Scan() {
			c.String(http.StatusOK, "%s", scanner.Text())
			return
		}
		c.Status(http.StatusGone)
		return
	})

	r.GET("/remove/:key", func(c *gin.Context) {
		key := c.Param("key")
		if _, ok := keys[key]; ok {
			delete(keys, key)
			c.Status(http.StatusOK)
			return
		}
		c.Status(http.StatusNotFound)
		return
	})

	log.Fatal(r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)))
}
