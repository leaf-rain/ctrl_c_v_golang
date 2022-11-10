package emoji

import (
	"fmt"
	emoji "github.com/Andrew-M-C/go.emoji"
	"strings"
)

// EmojiNameEncrypt 带emoji name处理 str：需要加密名称， encrypt:加密字符，preNum:前面预留字符数量，sufNum:后面预留字符数量
func EmojiNameEncrypt(str, encrypt string, preNum, sufNum int) string {
	i := 0
	var emojiMap = make(map[string]string, 0)
	var keySlice = make([]string, 0)
	final := emoji.ReplaceAllEmojiFunc(str, func(emoji string) string {
		var key = fmt.Sprintf("{{%d}}", i)
		emojiMap[key] = emoji
		keySlice = append(keySlice, key)
		i += 1
		return key
	})
	var result = make([]string, 0)
	var left, right int
	var pre string
	if len(keySlice) > 0 {
		for index, item := range keySlice {
			if index == 0 {
				left = 0
				right = strings.Index(final, item)
			} else {
				pre = keySlice[index-1]
				left = strings.Index(final, pre) + len(pre)
				right = strings.Index(final, item)
			}
			for _, strItem := range final[left:right] {
				result = append(result, string(strItem))
			}
			result = append(result, item)
			if index == len(keySlice)-1 { // 检查是否还有继续
				right += len(item)
				for _, endItem := range final[right:] {
					result = append(result, string(endItem))
				}
			}
		}
	} else {
		result = strings.Split(str, "")
	}
	for index := range result {
		if _, ok := emojiMap[result[index]]; ok {
			result[index] = emojiMap[result[index]]
		}
	}
	return charEncry(result, preNum, sufNum, encrypt)
}

func charEncry(name []string, preNum, sufNum int, encrypt string) string {
	var length = len(name)
	if length < preNum+sufNum {
		var right int
		if length >= preNum {
			right = preNum
		} else {
			right = length
		}
		return strings.Join(name[:right], "") + encrypt
	} else {
		return strings.Join(name[:preNum], "") + encrypt + strings.Join(name[length-sufNum:], "")
	}
}
