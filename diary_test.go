package diary

import (
	// "os"
	"testing"
)

func TestPass(t *testing.T) {
	t.Logf(string(generateHash("baek4531")))
	t.Log(passwordHash, "baek4531")
}

func _TestYM(t *testing.T) {
	t.Log(yearMonthConvert("2012-12"))
	t.Log(yearMonthConvert("2012-12~2013-12"))
	t.Log(yearMonthConvert("2012-12~2014-12"))
	t.Log(yearMonthConvert("2012-12~2015-12"))
}

func _TestDairy(t *testing.T) {
	d := getJustWrittenDiaryFromDirectory()
	prepareDirectoryToSaveEncryptedDiary(d)
	saveDiaryToDirectory(d)
}

func _TestOpenDiary(t *testing.T) {
	var diary []*Diary
	t.Log(openDiaryfromDirectory(&diary, "./2015/12"))
}
