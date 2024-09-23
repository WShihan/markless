package tool

import (
	"fmt"
	"html/template"
	"log/slog"
	"markless/model"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func FileOrPathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func GetBaseTemplate() *template.Template {
	return template.New("templates/temp.html")
}
func ExcutePath() string {
	excutePath, err := os.Executable()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return filepath.Dir(excutePath)
}

// 自定义函数，用于合并多个 User 对象的 FirstName 字段
func JoinTagNames(tags []model.Tag) string {
	firstNames := make([]string, len(tags))
	for i, user := range tags {
		firstNames[i] = user.Name
	}
	return strings.Join(firstNames, "&")
}

// 自增模板函数
func Increase(num int) int {
	return num + 1
}

// 自增模板函数
func Decrease(num int) int {
	return num - 1
}

func TimeFMT(t time.Time) string {
	day := t.Day()
	month := t.Month()
	year := t.Year()
	hour := t.Hour()
	min := t.Minute()

	return fmt.Sprintf("%d-%d-%d %d:%d", year, month, day, hour, min)
}

func ShortUID(length int) (uid string) {
	ubyte := []byte(uuid.New().String())
	idByte := new(big.Int).SetBytes(ubyte)
	rawID := idByte.String()
	start := rand.Intn(len(rawID) - length - 1)
	uid = rawID[start : length+start]
	return
}

func HashMessage(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidateHash(hashedPassword, password string) error {
	// 验证密码
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func DefaultLanguage() string {
	return "zh-CN"
}

func RandomN() int {
	return rand.Intn(100)
}
