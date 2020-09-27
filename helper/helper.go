package helper 


import (
	//"fmt"
	"net/url"
	"math/rand"
	//"log"
	//"time" 
	"errors" 
	"os"
	"encoding/json"
 )


type Config struct {
	PORT int
	DATA_SOURCE_NAME string
}

func (c *Config) Load() error{
	file, _ := os.Open("config.json")
	defer file.Close()
	err := json.NewDecoder(file).Decode(c)
	if err != nil {
		errors.New("config reading error")
	}

        return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
func RandSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}


func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

