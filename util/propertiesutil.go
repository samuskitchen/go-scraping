/**
* Utility that reads or traverses a .properties file
**/
package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/**
* Local structure of the utility to obtain or read the values of the properties
**/
type Properties struct {
	value map[string]string
	path  string
}

/**
* A properties file is created for GO reading
**/
func NewProperties() (result *Properties) {

	path, err := filepath.Abs("scraping.properties")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	result = &Properties{}
	if len(path) == 0 {
		fmt.Println("the route is incorrect.")
		return nil
	}

	result.path = path
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	result.value = getProperties(file)

	return
}

/**
* You get the properties of the .properties file
**/
func getProperties(file *os.File) map[string]string {
	reader := bufio.NewReader(file)
	result := make(map[string]string)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		s := string(line)
		s = strings.Replace(s, " ", "", -1)
		if len(s) == 0 || strings.Index(s, "#") == 0 {
			continue
		}
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			continue
		}
		result[kv[0]] = kv[1]
	}

	return result
}

/**
* Function that gets the value in string
**/
func (p Properties) GetString(key string) string {
	v, _ := p.value[key]
	return v
}

/**
* Function that gets the value in int
**/
func (p Properties) GetInteger(key string) int {
	v, _ := p.value[key]
	i, _ := strconv.Atoi(v)
	return i
}

/**
* Function that obtains the value in float
**/
func (p Properties) GetFloat(key string) float64 {
	v, _ := p.value[key]
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

/**
* Function that gets the value in bool
**/
func (p Properties) GetBool(key string) bool {
	v, _ := p.value[key]
	b, _ := strconv.ParseBool(v)
	return b
}

/**
* Function that is or inserts a value, to the position or value of the properties
**/
func (p *Properties) set(k, v string) {
	p.value[k] = v
}

/**
* Function traverses the properties to obtain the general information of the values
**/
func (p Properties) string() string {
	result := ""
	result += "value = ["
	for k, v := range p.value {
		result += k + "=" + v + ", "
	}
	result = string([]rune(result)[0 : len(result)-2])
	result += "],"
	result += "path = " + p.path
	return result
}