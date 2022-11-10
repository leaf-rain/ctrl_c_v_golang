package emoji

import (
	"fmt"
	"testing"
)

func TestEmojiNameEncrypt(t *testing.T) {
	fmt.Println(EmojiNameEncrypt("👩‍👩‍👦🇨🇳chinaNo1", "***", 5, 3))
}
