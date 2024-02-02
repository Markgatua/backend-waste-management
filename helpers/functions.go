package helpers

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
	"ttnmwastemanagementsystem/configs"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/models"
)

type Functions struct{}

func (functions Functions) TokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (functions Functions) GetRandString(length int) string {
	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b) // generates len(b) random bytes
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}

func (functions Functions) NumberTokenGenerator(numberOfDigits int) (int, error) {
	maxLimit := int64(int(math.Pow10(numberOfDigits)) - 1)
	lowLimit := int(math.Pow10(numberOfDigits - 1))

	randomNumber, err := rand.Int(rand.Reader, big.NewInt(maxLimit))
	if err != nil {
		return 0, err
	}
	randomNumberInt := int(randomNumber.Int64())

	// Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
	if randomNumberInt <= lowLimit {
		randomNumberInt += lowLimit
	}

	// Never likely to occur, kust for safe side.
	if randomNumberInt > int(maxLimit) {
		randomNumberInt = int(maxLimit)
	}
	return randomNumberInt, nil
}

func (functions Functions) FileToString(filePath string) string {
	b, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		return ""
	}

	str := string(b) // convert content to a 'string'
	return str
}

func (functions Functions) HashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func (functions Functions) GenerateToken(userId int64) (string, error) {

	tokenLifespan := configs.EnvConfigs.JWTExp
	secret := configs.EnvConfigs.JWTSecret

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}

func (functions Functions) TokenValid(c *gin.Context) error {
	secret := configs.EnvConfigs.JWTSecret

	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (functions Functions) ReplaceTemplateWithOrganizationInformation(template string) string {

	organization := models.TtnmOrganizationModel{}
	err := gen.REPO.DB.Get(&organization, "select name, logo_path, website_url from ttnm_organization")
	if err != nil {
		organization.Name = ""
		organization.LogoPath = ""
		organization.WebsiteUrl = ""
	}

	fmt.Print(organization.WebsiteUrl)
	fmt.Print(organization.LogoPath)
	fmt.Println(template)

	newTemplate := template
	newTemplate = strings.ReplaceAll(newTemplate, "{{.organization_website_url}}", organization.WebsiteUrl)
	newTemplate = strings.ReplaceAll(newTemplate, "{{.organization_logo_path}}", organization.LogoPath)
	newTemplate = strings.ReplaceAll(newTemplate, "{{.organization_name}}", organization.Name)

	// newTemplate := strings.ReplaceAll(template, "{{.organization_website_url}}", organization.WebsiteUrl)
	// newTemplate = strings.ReplaceAll(template, "{{.organization_logo_path}}", organization.LogoPath)
	// newTemplate = strings.ReplaceAll(template, "{{.organization_name}}", organization.Name)

	fmt.Println(newTemplate)
	return newTemplate
}

func (functions Functions) SelectScan(rows *sql.Rows) ([]map[string]interface{}, error) {
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	numColumns := len(columns)

	values := make([]interface{}, numColumns)
	for i := range values {
		values[i] = new(interface{})
	}

	var results []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		dest := make(map[string]interface{}, numColumns)
		for i, column := range columns {
			dest[column] = *(values[i].(*interface{}))
		}
		results = append(results, dest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (functions Functions) CurrentUserFromToken(c *gin.Context) (*models.User, error) {
	secret := configs.EnvConfigs.JWTSecret

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["id"]), 10, 32)
		if err != nil {
			return nil, err
		}

		user := models.User{}
		err = gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where id=?"), userID)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, errors.New("Error occured")
}
