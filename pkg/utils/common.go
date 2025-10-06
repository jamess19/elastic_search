package utils

import (
	"business/pkg/model"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func StrDelimitForSum(flt float64, currency string) string {
	str := strconv.FormatFloat(flt, 'f', 0, 64)

	pos := len(str) - 3
	for pos > 0 {
		str = str[:pos] + "." + str[pos:]
		pos = pos - 3
	}

	if currency != "" {
		return str + " " + currency
	}
	return str
}

func ParseIDFromUri(c *gin.Context) *uuid.UUID {
	tID := model.UriParse{}
	if err := c.ShouldBindUri(&tID); err != nil {
		_ = c.Error(err)
		return nil
	}
	if len(tID.ID) == 0 {
		_ = c.Error(fmt.Errorf("error: Empty when parse ID from URI"))
		return nil
	}
	if id, err := uuid.Parse(tID.ID[0]); err != nil {
		_ = c.Error(err)
		return nil
	} else {
		return &id
	}
}

func ParseStringIDFromUri(c *gin.Context) *string {
	tID := model.UriParse{}
	if err := c.ShouldBindUri(&tID); err != nil {
		_ = c.Error(err)
		return nil
	}
	if len(tID.ID) == 0 {
		_ = c.Error(fmt.Errorf("error: Empty when parse ID from URI"))
		return nil
	}
	return &tID.ID[0]
}

func ResizeImage(link string, w, h int) string {
	if link == "" || w == 0 || !strings.Contains(link, LINK_IMAGE_RESIZE) {
		return link
	}

	size := getSizeImage(w, h)

	env := "/finan-dev/"
	linkTemp := strings.Split(link, "/finan-dev/")
	if len(linkTemp) != 2 {
		linkTemp = strings.Split(link, "/finan/")
		env = "/finan/"
	}

	if len(linkTemp) == 2 {
		url := linkTemp[0] + "/v2/" + size + env + linkTemp[1]
		return strings.ReplaceAll(url, " ", "%20")
	}
	return strings.ReplaceAll(link, " ", "%20")
}

func getSizeImage(w, h int) string {
	if h == 0 {
		return "w" + strconv.Itoa(w)
	}
	return strconv.Itoa(w) + "x" + strconv.Itoa(h)
}

func ConvertVNPhoneFormat(phone string) string {
	if phone != "" {
		if strings.HasPrefix(phone, "84") {
			phone = "+" + phone
		}
		if strings.HasPrefix(phone, "0") {
			phone = "+84" + phone[1:]
		}
	}
	return phone
}

func ValidPhoneFormat(phone string) bool {
	if phone == "" {
		return false
	}
	if len(phone) == 13 {
		return true
	}
	internationalPhone := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	vietnamPhone := regexp.MustCompile(`((09|03|07|08|05)+([0-9]{8})\b)`)
	if !vietnamPhone.MatchString(phone) {
		if !internationalPhone.MatchString(phone) {
			return false
		}
	}
	return true
}

func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func RevertBeginPhone(phone string) string {
	if phone != "" {
		if strings.HasPrefix(phone, "+84") {
			phone = "0" + phone[3:]
		}
	}
	return phone
}

func CurrentFunctionName(level int) string {
	pc, _, _, _ := runtime.Caller(1 + level)
	strArr := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	if len(strArr) > 0 {
		return strArr[len(strArr)-1]
	}
	return ""
}

func GetCurrentCaller(caller interface{}, level int) string {
	strArr := strings.Split(reflect.TypeOf(caller).String(), ".")
	if caller != nil && len(strArr) > 0 {
		return fmt.Sprintf("%s.%s", strArr[len(strArr)-1], CurrentFunctionName(1+level))
	}
	return CurrentFunctionName(1)
}

func CheckRequireValid(ob interface{}) error {
	validator := validation.Validation{RequiredFirst: true}
	passed, err := validator.Valid(ob)
	if err != nil {
		return err
	}
	if !passed {
		var err string
		for _, e := range validator.Errors {
			err += fmt.Sprintf("[%s: %s] ", e.Field, e.Message)
		}
		return fmt.Errorf("%s", err)
	}
	return nil
}
