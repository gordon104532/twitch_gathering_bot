package TwitchBot

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/app/ErrorHandle"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ReadPointFile() (result string) {
	filename := "gatTotalPoint.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			ErrorHandle.Error.Println("ERROR", "readPointFile: 讀取分數檔案錯誤, "+err.Error())
		}
	}()

	b, err := ioutil.ReadAll(file)
	return string(b)
}

// 輸入名次列出前幾名與累計分數
func RankByPoint(target int) (rank string) {
	// 讀分數
	source := ReadPointFile()

	if len(source) == 0 {
		return "Empty"
	}

	// [暱稱]:分數
	result := make(map[string]int, 0)

	// 字串分割
	splitStr := strings.Split(source, ":")

	for i, str := range splitStr {
		if strings.Contains(str, "暱稱") {
			// 只切出暱稱
			aliasRaw := strings.Split(splitStr[i+1], ",")

			// 找出分數
			pointRaw := strings.Split(splitStr[i+2], ",")

			point, err := strconv.Atoi(strings.TrimSpace(pointRaw[0]))
			if err != nil {
				ErrorHandle.Error.Println("ERROR", "strconv.Atoi err:, "+err.Error())
			}
			result[aliasRaw[0]] += point
		}
	}

	// 排序分數
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range result {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var printRank int
	if len(result) < target {
		printRank = len(result)
	} else {
		printRank = target
	}

	for i := 0; i < printRank; i++ {
		rank += fmt.Sprintf("Rank%d: %s, %d\n", i+1, ss[i].Key, ss[i].Value)
	}

	fmt.Print(rank)
	fmt.Println("-----")
	return
}
