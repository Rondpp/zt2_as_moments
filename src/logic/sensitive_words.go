package logic

import (
	"bufio"
	"conf"
	"fmt"
	"github.com/eachain/aca"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

var sensitive_words *aca.ACA

func init() {
	sensitive_words = aca.New()
	fmt.Println(conf.GetCfg())
	f, err := os.Open("conf/sensitive_words.txt") //conf.GetCfg().SensitiveWordsFile)
	if err != nil {
		fmt.Println(err, conf.GetCfg().SensitiveWordsFile)
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		if len(line) >= 2 {
			sensitive_words.Add(line[0 : len(line)-1])
		}
	}
	sensitive_words.Build()
}

func SensitiveWordsReplace(words *string) bool {
	word_arr := sensitive_words.Find(*words)

	for _, v := range word_arr {
		var replace_str string
		i := 0
		for i < utf8.RuneCountInString(v) {
			i = i + 1
			replace_str += "*"
		}
		*words = strings.Replace(*words, v, replace_str, -1)
	}
	return len(word_arr) > 0
}
