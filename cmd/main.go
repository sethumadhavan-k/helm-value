package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

func getKey(baseKey string, key string) string {
	if baseKey != "" {
		return baseKey + "." + key
	}
	return key
}

func stringConcat(str string, value interface{}) string {
	return str + fmt.Sprintf("%v", value) + " "
}

func convertTypeAndAppend(value interface{}, valueType reflect.Kind, key string, keyArray *[][]string) {
	switch valueType {
	case reflect.Int:
		*keyArray = append(*keyArray, []string{key, strconv.Itoa(value.(int))})
	case reflect.Float64:
		*keyArray = append(*keyArray, []string{key, fmt.Sprintf("%v", value)})
	case reflect.Bool:
		*keyArray = append(*keyArray, []string{key, strconv.FormatBool(value.(bool))})
	case reflect.String:
		*keyArray = append(*keyArray, []string{key, value.(string)})
	}
}

func getHelmValues(baseKey string, data map[string]interface{}, keys *[][]string) {
	for key, value := range data {
		valueType := reflect.ValueOf(value).Kind()
		switch valueType {
		case reflect.Slice:
			sv := "["
			sliceValue := value.([]interface{})
			for _, svalue := range sliceValue {
				sv = stringConcat(sv, svalue)
			}
			sv = sv + "]"
			*keys = append(*keys, []string{getKey(baseKey, key), sv})
		case reflect.Map:
			getHelmValues(getKey(baseKey, key), value.(map[string]interface{}), keys)
		default:
			convertTypeAndAppend(value, valueType, getKey(baseKey, key), keys)
		}
	}
}

func rendeTable(row [][]string, search string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Propery", "Default Value"})
	sort.SliceStable(row, func(p, q int) bool {
		return row[p][0] < row[q][0]
	})
	for _, v := range row {
		if search != "" && !strings.Contains(v[0], search) {
			continue
		}
		table.Append(v)

	}
	table.Render()
}

func App(repo string, search string) {
	fmt.Println(search)
	stdout, err := exec.Command("helm", "show", "values", repo).Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var helmYAML map[string]interface{}

	err2 := yaml.Unmarshal([]byte(string(stdout)), &helmYAML)
	if err2 != nil {
		fmt.Println("error: %v", err2)
	}
	var helmValueSlice [][]string
	getHelmValues("", helmYAML, &helmValueSlice)
	rendeTable(helmValueSlice, search)
}
