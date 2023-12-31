package controller

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spurtcms/spurtcms-content/lms"
	"github.com/spurtcms/spurtcms-core/auth"
	"github.com/spurtcms/spurtcms-core/member"
	"gorm.io/gorm"

	spurtcore "github.com/spurtcms/spurtcms-core"
)

type key string

const ContextKey key = "ginContext"

var IST, _ = time.LoadLocation("Asia/Kolkata")

var ProfileImagePath = "Uploads/ProfileImages/"

var(
	Mem member.MemberAuth
	Auth *auth.Authorization
	Flag bool
	Pg lms.MemberPage
	Sp lms.MemberSpace
)

func GetAuthorization(token string,db *gorm.DB)(*auth.Authorization) {

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: token, Secret: os.Getenv("JWT_SECRET")})

	return &auth

}

func GetAuthorizationWithoutToken(db *gorm.DB)(*auth.Authorization){

	auth := spurtcore.NewInstance(&auth.Option{DB: db, Token: "", Secret: ""})

	return &auth
}

func StoreImageBase64ToLocal(imageData,storagePath,storingName string) (string,string,error) {

	extEndIndex := strings.Index(imageData, ";base64,")

	base64data := imageData[strings.IndexByte(imageData, ',')+1:]

	var ext = imageData[11:extEndIndex]

	randomNum := strconv.Itoa(rand.Intn(900000) + 100000)

	imageName := storingName +"-"+ randomNum + "." + ext

	err := os.MkdirAll(storagePath, 0755)

	if err!=nil{

		log.Println(err)

		return "","",err
	}

	storageDestination := storagePath + imageName
	
	decode, err := base64.StdEncoding.DecodeString(base64data)

	if err != nil {

		log.Println(err)

		return "","",err
	}

	file, err := os.Create(storageDestination)

	if err != nil {

		log.Println(err)

		return "","",err

	}
	if _, err := file.Write(decode); err != nil {

		log.Println(err)

		return "","",err

	}

	return imageName,storageDestination,nil
}

func SendEmail(userEmail,sub,body string ,channel chan error){

	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	to := []string{
		userEmail,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject:" + sub + "\n"

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)

	if err != nil {

		log.Println(err)

		channel <- err

		return
	}

	channel <- nil

}


