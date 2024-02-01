package controllers

import (
	"ttnmwastemanagementsystem/configs"
	"ttnmwastemanagementsystem/gen"

	// "ttnmwastemanagementsystem/firebaseapp"
	"database/sql"
	"errors"
	_ "errors"
	"fmt"
	"net/http"
	"time"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"
	"ttnmwastemanagementsystem/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	// "gopkg.in/guregu/null.v3"
	// "gopkg.in/guregu/null.v3"
)

const TAG string = "controllers/Auth"

type EmailResetPasswordParam struct {
	Email string `json:"email"  binding:"required"`
}

type SendVerificationMailParam struct {
	Email string `json:"email"`
}

type VerifyOTPCodePhoneParam struct {
	OTPCode     string `json:"otp_code" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	CallingCode string `json:"calling_code" binding:"required"`
}

type SendOTPCodePhoneParam struct {
	Phone       string `json:"phone" binding:"required"`
	CallingCode string `json:"calling_code" binding:"required"`
}

type LoginPhoneParam struct {
	Phone       string `json:"phone" binding:"required"`
	CallingCode string `json:"calling_code" binding:"required"`
	Pin         string `json:"pin" binding:"required"`
}

type ForgotPinSendOTPPhoneParam struct {
	Phone       string `json:"phone" binding:"required"`
	CallingCode string `json:"calling_code" binding:"required"`
}

type ForgotPinVerifyOTPPhoneParam struct {
	OTPCode     string `json:"otp_code" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	CallingCode string `json:"calling_code" binding:"required"`
}

type ForgotPinEnterNewPinParam struct {
	RecoveryCode string `json:"recovery_code" binding:"required"`
	Pin          string `json:"pin" binding:"required"`
}

type RegisterUserEmailParam struct {
	Email         string `json:"email" binding:"required"`
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Password      string `json:"password" binding:"required"`
	UserType      int32  `json:"user_type" binding:"required"`
	RoleId        int32  `json:"role_id" binding:"required"`
	UserCompanyId int32  `json:"user_company_id"`
}

type UpdatePasswordParam struct {
	Password string `json:"password" binding:"required"`
	UserID   int32  `json:"user_id" binding:"required"`
}

type EditUserEmailParam struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	RoleId    int32  `json:"role_id" binding:"required"`
	UserID    int32  `json:"user_id" binding:"required"`
}

type ResetPasswordApiParams struct {
	Email    string `json:"email" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordPhoneApiParams struct {
	Phone       string `json:"phone" binding:"required"`
	CallingCode string `json:"calling_code" binding:"required"`
	OTPCode     string `json:"otp_code" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type LoginUserEmailParam struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterPhoneUpdateUserDetailsParam struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Pin       string `json:"pin" binding:"required"`
}

type PhoneVerificationToken struct {
	ID          int64     `db:"id" json:"id"`
	Token       string    `db:"token" json:"token"`
	CallingCode string    `db:"calling_code" json:"calling_code"`
	Phone       string    `db:"phone" json:"phone"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type AuthController struct{}

func IsVerificationTokenValid(token string) (bool, error) {
	tokenExpirationTime := configs.EnvConfigs.AcoountVerificationTokenExpirationTime
	//	if !tokenExpirationTime.Valid {/
	//		return false, errors.New("Error fetching token expiration time")
	//	}
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select confirmation_token,confirmation_sent_at from users where confirmation_token=?"), token)
	if err != nil {
		return false, err
	}

	now := carbon.Now()

	confirmationSentAtTime := carbon.CreateFromDateTime(user.ConfirmationSentAt.Time.Year(), int(user.ConfirmationSentAt.Time.Month()), user.ConfirmationSentAt.Time.Day(), user.ConfirmationSentAt.Time.Hour(), user.ConfirmationSentAt.Time.Minute(), user.ConfirmationSentAt.Time.Second())
	diff := now.DiffAbsInHours(confirmationSentAtTime)

	logger.Log(TAG, fmt.Sprint("Time now :: ", now), logger.LOG_LEVEL_INFO)
	logger.Log(TAG, fmt.Sprint("Confrimation sent at :: ", confirmationSentAtTime), logger.LOG_LEVEL_INFO)
	logger.Log(TAG, fmt.Sprint("Difference in hours ::", diff), logger.LOG_LEVEL_INFO)
	logger.Log(TAG, fmt.Sprint("Token expiration in hours ::", tokenExpirationTime), logger.LOG_LEVEL_INFO)

	if diff >= int64(tokenExpirationTime) {
		return false, nil
	} else {
		return true, nil
	}
}

func IsRecoveryTokenValid(token string) (bool, error) {

	tokenExpirationTime := configs.EnvConfigs.AccountRecoveryTokenExpirationTime
	//	if !tokenExpirationTime.Valid {
	//		return false, errors.New("Error fetching token expiration time")
	//	}
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select recovery_token,recovery_sent_at from users where recovery_token=?"), token)
	if err != nil {
		return false, err
	}
	now := carbon.Now()
	recoverySentAtTime := carbon.CreateFromDateTime(user.RecoverySentAt.Time.Year(), int(user.RecoverySentAt.Time.Month()), user.RecoverySentAt.Time.Day(), user.RecoverySentAt.Time.Hour(), user.RecoverySentAt.Time.Minute(), user.RecoverySentAt.Time.Second())
	diff := now.DiffAbsInHours(recoverySentAtTime)

	logger.Log(TAG, fmt.Sprint("Time now :: ", now), logger.LOG_LEVEL_INFO)
	logger.Log(TAG, fmt.Sprint("Recovery sent at :: ", recoverySentAtTime), logger.LOG_LEVEL_INFO)
	logger.Log(TAG, fmt.Sprint("Difference in hours ::", diff), logger.LOG_LEVEL_INFO)
	logger.Log(TAG, fmt.Sprint("Token expiration in hours ::", tokenExpirationTime), logger.LOG_LEVEL_INFO)

	if diff >= int64(tokenExpirationTime) {
		return false, errors.New("Token has expired")
	} else {
		return true, nil
	}
}

func (auth AuthController) SendVerificationMail(context *gin.Context) {
	var sendVerificationMailParam SendVerificationMailParam
	err := context.ShouldBindJSON(&sendVerificationMailParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := GetEmailUser(sendVerificationMailParam.Email)

	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User does not exist",
		})
		return
	}

	mail := helpers.Mail{}
	sent := mail.SendMailVerification(user.Email.String, fmt.Sprint(user.FirstName.String, " ", user.LastName.String))
	if sent {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Email verification link sent to your email",
		})
	} else {
		context.JSON(http.StatusExpectationFailed, gin.H{
			"error":   true,
			"message": "Error sending verification link to email",
		})
	}
}

func GetPhoneUser(callingCode string, phone string) (*models.User, error) {
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where calling_code=? and phone=?"), callingCode, phone)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID int64) (*models.User, error) {
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where id=?"), userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetEmailUser(email string) (*models.User, error) {
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where email=?"), email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUserExistingOnUpdate(email string, userID int32) (bool, error) {
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where email=? and id !=?"), email, userID)
	if err != nil {
		return false, err
	}
	fmt.Println(user)
	return user.ID.Valid, nil
}

func GetUserFromRecoveryTokenWithNoValidation(recoveryToken string) (*models.User, error) {
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where recovery_token=?"), recoveryToken)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserFromRecoveryTokenWithValidationVerification(recoveryToken string, callingCode string, phone string) (*models.User, error) {
	user := models.User{}
	err := gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select * from users where recovery_token=? and calling_code=? and phone=?"), recoveryToken, callingCode, phone)
	if err != nil {
		return nil, err
	}

	valid, err := IsRecoveryTokenValid(recoveryToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("Invalid token")
	}

	return &user, nil
}

// If the otp code is verified the phone is registered and a wallet is created for the user
func (auth AuthController) RegisterVerifyOTPCodePhoneAndCreateWallet(context *gin.Context) {
	var verifyOTPCodePhoneParam VerifyOTPCodePhoneParam
	err := context.ShouldBindJSON(&verifyOTPCodePhoneParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	phoneVerificationToken := PhoneVerificationToken{}
	err = gen.REPO.DB.Get(&phoneVerificationToken, gen.REPO.DB.Rebind("select * from phone_verification_token where token=?"), verifyOTPCodePhoneParam.OTPCode)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error registering user.",
		})
		return
	}

	user, _ := GetPhoneUser(phoneVerificationToken.CallingCode, phoneVerificationToken.Phone)
	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User already exists",
		})
		return
	}
	_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (calling_code,phone,user_type,created_at,updated_at,provider,phone_confirmed_at) VALUES (:calling_code,:phone,:user_type,:created_at,:updated_at,:provider,:phone_confirmed_at)`,
		map[string]interface{}{
			"user_type":          1,
			"calling_code":       phoneVerificationToken.CallingCode,
			"phone":              phoneVerificationToken.Phone,
			"provider":           "phone",
			"created_at":         time.Now(),
			"updated_at":         time.Now(),
			"phone_confirmed_at": time.Now(),
		})

	//fmt.Println(err.Error())

	user, err = GetPhoneUser(phoneVerificationToken.CallingCode, phoneVerificationToken.Phone)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error registering user",
		})
		return
	}

	// _, err = gen.REPO.DB.NamedExec(`INSERT INTO user_wallet (user_id,wallet_ref,amount,created_at) VALUES (:user_id,:wallet_ref,:amount,:created_at)`,
	// 	map[string]interface{}{
	// 		"user_id":    user.ID.Int64,
	// 		"wallet_ref": helpers.Functions{}.WalletRef(),
	// 		"amount":     0,
	// 		"created_at": time.Now(),
	// 	})

	//gen.REPO.DB.NamedExec(`delete from phone_verification_token where calling_code=:calling_code and phone=:phone`,
	//	map[string]interface{}{
	//		"calling_code": phoneVerificationToken.CallingCode,
	//		"phone":        phoneVerificationToken.Phone,
	//	})

	token, err := helpers.Functions{}.GenerateToken(user.ID.Int64)
	// wallet := helpers.Wallet{}.GetWalletForUser(user.ID.Int64)

	// go func() {
	// 	firebaseapp.UpdateCurrentAmountForUser(fmt.Sprint(phoneVerificationToken.CallingCode, phoneVerificationToken.Phone), 0)
	// }()

	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user,
		// "wallet_ref": wallet.WalletRef.String,
		"token": token,
	})
}

// / If the otp code is verified the phone is registered
// /
func (auth AuthController) RegisterVerifyOTPCodePhone(context *gin.Context) {
	var verifyOTPCodePhoneParam VerifyOTPCodePhoneParam
	err := context.ShouldBindJSON(&verifyOTPCodePhoneParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	phoneVerificationToken := PhoneVerificationToken{}
	err = gen.REPO.DB.Get(&phoneVerificationToken, gen.REPO.DB.Rebind("select * from phone_verification_token where token=? and calling_code=? and phone=?"), verifyOTPCodePhoneParam.OTPCode, verifyOTPCodePhoneParam.CallingCode, verifyOTPCodePhoneParam.Phone)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error registering user.",
		})
		return
	}

	user, _ := GetPhoneUser(phoneVerificationToken.CallingCode, phoneVerificationToken.Phone)
	// if user != nil {
	// 	context.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"error":   true,
	// 		"message": "User already exists",
	// 	})
	// 	return
	// }

	if user == nil {
		_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (calling_code,phone,user_type,created_at,updated_at,provider,phone_confirmed_at) VALUES (:calling_code,:phone,:user_type,:created_at,:updated_at,:provider,:phone_confirmed_at)`,
			map[string]interface{}{
				"user_type":          1,
				"calling_code":       phoneVerificationToken.CallingCode,
				"phone":              phoneVerificationToken.Phone,
				"provider":           "phone",
				"created_at":         time.Now(),
				"updated_at":         time.Now(),
				"phone_confirmed_at": time.Now(),
			})
	}

	//fmt.Println(err.Error())

	user, err = GetPhoneUser(phoneVerificationToken.CallingCode, phoneVerificationToken.Phone)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error registering user",
		})
		return
	}

	//gen.REPO.DB.NamedExec(`delete from phone_verification_token where calling_code=:calling_code and phone=:phone`,
	//	map[string]interface{}{
	//		"calling_code": phoneVerificationToken.CallingCode,
	//		"phone":        phoneVerificationToken.Phone,
	//	})

	token, err := helpers.Functions{}.GenerateToken(user.ID.Int64)
	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user,
		"token": token,
	})
}

func (auth AuthController) GetUser(context *gin.Context) {
	currentUser, err := helpers.Functions{}.CurrentUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Error getting user",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  currentUser,
	})
	return

}

func (auth AuthController) ForgotPinEnterNewPin(context *gin.Context) {
	var forgotPinEnterNewPinParam ForgotPinEnterNewPinParam
	err := context.ShouldBindJSON(&forgotPinEnterNewPinParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetUserFromRecoveryTokenWithNoValidation(forgotPinEnterNewPinParam.RecoveryCode)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error reseting pin",
		})
		return
	}
	_, err = gen.REPO.DB.NamedExec(`UPDATE users SET password=:password where id=:id`, map[string]interface{}{
		"password": helpers.Functions{}.HashPassword(forgotPinEnterNewPinParam.Pin),
		"id":       user.ID,
	})
	if err != nil {
		logger.Log(TAG, fmt.Sprint("Query error:", err), logger.LOG_LEVEL_ERROR)
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Pin reset success",
		})
	}
}

func (auth AuthController) ForgotPinVerifyOTPPhone(context *gin.Context) {
	var forgotPinVerifyOTPPhoneParam ForgotPinVerifyOTPPhoneParam
	err := context.ShouldBindJSON(&forgotPinVerifyOTPPhoneParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetUserFromRecoveryTokenWithValidationVerification(forgotPinVerifyOTPPhoneParam.OTPCode, forgotPinVerifyOTPPhoneParam.CallingCode, forgotPinVerifyOTPPhoneParam.Phone)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   true,
		"message": "OTP code verified successfully",
	})
}

func (auth AuthController) ForgotPinSendOTPPhone(context *gin.Context) {
	var forgotPinSendOTPPhoneParam ForgotPinSendOTPPhoneParam
	err := context.ShouldBindJSON(&forgotPinSendOTPPhoneParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, _ := GetPhoneUser(forgotPinSendOTPPhoneParam.CallingCode, forgotPinSendOTPPhoneParam.Phone)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given phone doesnt exist",
		})
		return
	}

	otpCode, _ := helpers.Functions{}.NumberTokenGenerator(4)

	_, err = gen.REPO.DB.NamedExec(`UPDATE users set recovery_token=:token,recovery_sent_at=:created_at where id=:id`,
		map[string]interface{}{
			"id":         user.ID,
			"token":      otpCode,
			"created_at": time.Now(),
		})

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error sending pin reset code",
		})
		return
	}

	phone := forgotPinSendOTPPhoneParam.CallingCode + forgotPinSendOTPPhoneParam.Phone
	sms := helpers.SMS{}
	err = sms.SendSMS([]string{phone}, fmt.Sprint("Welcome to TakaTaka Ni Mali, your pin reset code: ", otpCode))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error sending pin reset code",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Pin reset code sent to your phone",
		})
	}

}

func (auth AuthController) LoginPhone(context *gin.Context) {
	var loginPhoneParam LoginPhoneParam
	err := context.ShouldBindJSON(&loginPhoneParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, _ := GetPhoneUser(loginPhoneParam.CallingCode, loginPhoneParam.Phone)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given phone doesnt exist",
		})
		return
	}

	functions := helpers.Functions{}
	hashedPassword := functions.HashPassword(loginPhoneParam.Pin)
	storedPassword := user.Password

	if hashedPassword != storedPassword.String {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Incorrect credentials",
		})
		return
	}

	token, err := helpers.Functions{}.GenerateToken(user.ID.Int64)
	// wallet := helpers.Wallet{}.GetWalletForUser(user.ID.Int64)

	context.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  user,
		// "wallet_ref": wallet.WalletRef.String,
		"token": token,
	})

}

func (auth AuthController) RegisterPhoneUpdateUserDetails(context *gin.Context) {
	var registerPhoneUpdateUserDetailsParam RegisterPhoneUpdateUserDetailsParam
	err := context.ShouldBindJSON(&registerPhoneUpdateUserDetailsParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	currentUser, err := helpers.Functions{}.CurrentUserFromToken(context)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Error getting current user",
		})
		return
	}

	_, err = gen.REPO.DB.NamedExec(`UPDATE users set first_name=:first_name,last_name=:last_name,password=:password where id=:id`,
		map[string]interface{}{
			"id":         currentUser.ID,
			"first_name": registerPhoneUpdateUserDetailsParam.FirstName,
			"last_name":  registerPhoneUpdateUserDetailsParam.LastName,
			"password":   helpers.Functions{}.HashPassword(registerPhoneUpdateUserDetailsParam.Pin),
		})

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Error saving user info",
		})
	} else {
		user, _ := GetUserByID(currentUser.ID.Int64)

		token, _ := helpers.Functions{}.GenerateToken(user.ID.Int64)
		// wallet := helpers.Wallet{}.GetWalletForUser(user.ID.Int64)

		context.JSON(http.StatusOK, gin.H{
			"error": false,
			"user":  user,
			// "wallet_ref": wallet.WalletRef.String,
			"token": token,
		})

	}

}

func (auth AuthController) RegisterPhoneSendOTPCode(context *gin.Context) {
	var sendOTPCodePhoneParam SendOTPCodePhoneParam
	err := context.ShouldBindJSON(&sendOTPCodePhoneParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, _ := GetPhoneUser(sendOTPCodePhoneParam.CallingCode, sendOTPCodePhoneParam.Phone)

	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given phone already exists",
		})
		return
	}

	otpCode, _ := helpers.Functions{}.NumberTokenGenerator(4)
	phone := sendOTPCodePhoneParam.CallingCode + sendOTPCodePhoneParam.Phone

	phoneVerificationToken := PhoneVerificationToken{}
	err = gen.REPO.DB.Get(&phoneVerificationToken, gen.REPO.DB.Rebind("select * from phone_verification_token where calling_code=? and phone=?"), sendOTPCodePhoneParam.CallingCode, sendOTPCodePhoneParam.Phone)
	//err if record does not exist

	if err != nil && err != sql.ErrNoRows {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error sending OTP code..",
		})
		return
	}
	if err != nil {
		//not found do insertion
		gen.REPO.DB.NamedExec(`INSERT INTO phone_verification_token (token,calling_code,phone,created_at) values(:token,:calling_code,:phone,:created_at)`,
			map[string]interface{}{
				"calling_code": sendOTPCodePhoneParam.CallingCode,
				"phone":        sendOTPCodePhoneParam.Phone,
				"token":        otpCode,
				"created_at":   time.Now(),
			})
		//fmt.Println(err.Error());
	} else {
		//do update
		gen.REPO.DB.NamedExec(`UPDATE phone_verification_token set token=:token,created_at=:created_at where id=:id`,
			map[string]interface{}{
				"id":         phoneVerificationToken.ID,
				"token":      otpCode,
				"created_at": time.Now(),
			})
	}

	sms := helpers.SMS{}
	err = sms.SendSMS([]string{phone}, fmt.Sprint("Welcome to TakaTaka Ni Mali, enter the OTP code : ", otpCode))
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error sending OTP code",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "OTP code sent to your phone",
		})
	}
}

func (auth AuthController) LoginEmail(context *gin.Context) {
	var loginUserEmailParam LoginUserEmailParam
	err := context.ShouldBindJSON(&loginUserEmailParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetEmailUser(loginUserEmailParam.Email)
	if user == nil {
		fmt.Print(err.Error())
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address does not exist",
		})
		return
	}

	functions := helpers.Functions{}
	hashedPassword := functions.HashPassword(loginUserEmailParam.Password)
	storedPassword := user.Password

	if hashedPassword != storedPassword.String {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Incorrect credentials",
		})
		return
	}

	if user.ConfirmedAt.Time.IsZero() {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Confirm Your Email",
		})
		return
	}

	token, err := helpers.Functions{}.GenerateToken(user.ID.Int64)
	permissionsForRole, _ := GetPermissionsForRole(int32(user.RoleId.Int64))
	actionList := GetActionsFromPermissions(permissionsForRole)

	context.JSON(http.StatusOK, gin.H{
		"error":       false,
		"user":        user,
		"permissions": actionList,
		"token":       token,
	})
}
func (auth AuthController) EditUserEmail(context *gin.Context) {
	organization, _ := context.Params.Get("organization")
	var editUserParam EditUserEmailParam
	err := context.ShouldBindJSON(&editUserParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetEmailUser(editUserParam.Email)
	if user != nil {
		if user.ID.Int64 != int64(editUserParam.UserID) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Another user has the same email",
			})
			return
		}
	}
	fmt.Println(organization)
	_, err = gen.REPO.DB.NamedExec(`update users set first_name=:first_name,last_name=:last_name,email=:email,role_id=:role_id where id=:user_id`,
		map[string]interface{}{
			"email":      editUserParam.Email,
			"first_name": editUserParam.FirstName,
			"last_name":  editUserParam.LastName,
			"role_id":    editUserParam.RoleId,
			"user_id":    editUserParam.UserID,
			"updated_at": time.Now(),
		})
	context.JSON(http.StatusUnprocessableEntity, gin.H{
		"error":   false,
		"message": "User updated successfully",
	})
	return
}
func (auth AuthController) UpdateUserPassword(context *gin.Context) {
	var param UpdatePasswordParam
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if len(param.Password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Password length should be greater than or equal to 6",
		})
		return
	}

	_, err = gen.REPO.DB.NamedExec(`update users set password=:password where id=:user_id`,
		map[string]interface{}{
			"password": helpers.Functions{}.HashPassword(param.Password),
			"user_id":  param.UserID,
		})

	if err != nil {
		fmt.Print(err.Error())
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error updating password",
		})
		return
	} else {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   false,
			"message": "Password updated successfully",
		})
		return
	}
}
func (auth AuthController) UpdateAggregatorUser(context *gin.Context) {
	type Params struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password"`
		CompanyID int32  `json:"institution_id" binding:"required"`
		RoleID    int32  `json:"role_id" binding:"required"`
		UserID    int32  `json:"user_id" binding:"required"`
		IsActive  *bool  `json:"is_active" binding:"required"`
	}

	var param Params
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	exists, err := IsUserExistingOnUpdate(param.Email, param.UserID)
	if err != nil && err != sql.ErrNoRows {
		logger.Log("AuthController", fmt.Sprint("Error adding user :: ", err.Error()), logger.LOG_LEVEL_ERROR)
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error adding user",
		})
		return
	}

	company, err := gen.REPO.GetCompany(context, param.CompanyID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if company.CompanyType != 2 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Selected aggregator institution",
		})
		return
	}

	if exists {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address already exists",
		})
		return
	}

	if param.Password != "" {
		_, err = gen.REPO.DB.NamedExec(`UPDATE users set first_name=:first_name,last_name=:last_name,email=:email,user_company_id=:user_company_id,role_id=:role_id,is_active=:is_active ,password=:password  where id=:id`,
			map[string]interface{}{
				"first_name":      param.FirstName,
				"last_name":       param.LastName,
				"email":           param.Email,
				"user_company_id": param.CompanyID,
				"role_id":         param.RoleID,
				"is_active":       param.IsActive,
				"id":              param.UserID,
				"password":        helpers.Functions{}.HashPassword(param.Password),
			})
		if err != nil && err != sql.ErrNoRows {
			logger.Log("AuthController [password set]", fmt.Sprint("Error adding user :: ", err), logger.LOG_LEVEL_ERROR)
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Error adding user",
			})
			return
		}
	} else {
		_, err = gen.REPO.DB.NamedExec(`UPDATE users set first_name=:first_name,last_name=:last_name,email=:email,user_company_id=:user_company_id,role_id=:role_id,is_active=:is_active where id=:id`,
			map[string]interface{}{
				"first_name":      param.FirstName,
				"last_name":       param.LastName,
				"email":           param.Email,
				"user_company_id": param.CompanyID,
				"role_id":         param.RoleID,
				"is_active":       param.IsActive,
				"id":              param.UserID,
			})

		if err != nil && err != sql.ErrNoRows {
			logger.Log("AuthController[password not set]", fmt.Sprint("Error adding user :: ", err), logger.LOG_LEVEL_ERROR)
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Error updating user params",
			})
			return
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "User updated successfully",
	})
}
func (auth AuthController) AddAggregatorUser(context *gin.Context) {
	type Params struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
		CompanyID int32  `json:"institution_id" binding:"required"`
		RoleID    int32  `json:"role_id" binding:"required"`
	}
	var param Params
	err := context.ShouldBindJSON(&param)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	company, err := gen.REPO.GetCompany(context, param.CompanyID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	if company.CompanyType != 2 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Selected aggregator institution",
		})
		return
	}

	user, err := GetEmailUser(param.Email)
	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address already exists",
		})
		return
	}
	if len(param.Password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Password length should be greater than or equal to 6",
		})
		return
	}

	_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,provider,role_id,user_company_id,created_at,updated_at,password,confirmed_at) VALUES (:email,:first_name,:last_name,:provider,:role_id,:user_company_id,:created_at,:updated_at,:password,:confirmed_at)`,
		map[string]interface{}{
			"email":           param.Email,
			"first_name":      param.FirstName,
			"last_name":       param.LastName,
			"role_id":         param.RoleID,
			"provider":        "email",
			"user_company_id": param.CompanyID,
			"password":        helpers.Functions{}.HashPassword(param.Password),
			"created_at":      time.Now(),
			"updated_at":      time.Now(),
			"confirmed_at":    time.Now(),
		})

	if err != nil && err != sql.ErrNoRows {
		logger.Log("AuthController", fmt.Sprint("Error adding user :: ", err.Error()), logger.LOG_LEVEL_ERROR)
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error adding user",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Aggregator user created successfully",
	})
}

func (auth AuthController) RegisterUserEmail(context *gin.Context) {
	organization, _ := context.Params.Get("organization")
	var registerUserEmailParam RegisterUserEmailParam
	err := context.ShouldBindJSON(&registerUserEmailParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetEmailUser(registerUserEmailParam.Email)
	if user != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address already exists",
		})
		return
	}

	if len(registerUserEmailParam.Password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Password length should be greater than or equal to 6",
		})
		return
	}

	if registerUserEmailParam.UserCompanyId != 0 {
		_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,user_type,provider,role_id,user_company_id,created_at,updated_at,password) VALUES (:email,:first_name,:last_name,:user_type,:provider,:role_id,:user_company_id,:created_at,:updated_at,:password)`,
			map[string]interface{}{
				"email":           registerUserEmailParam.Email,
				"first_name":      registerUserEmailParam.FirstName,
				"last_name":       registerUserEmailParam.LastName,
				"user_type":       registerUserEmailParam.UserType,
				"role_id":         registerUserEmailParam.RoleId,
				"provider":        "email",
				"user_company_id": registerUserEmailParam.UserCompanyId,
				"password":        helpers.Functions{}.HashPassword(registerUserEmailParam.Password),
				"created_at":      time.Now(),
				"updated_at":      time.Now(),
			})

	} else {

		if organization == "main_organization" {
			_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,user_type,provider,role_id,created_at,updated_at,password,is_main_organization_user,confirmed_at) VALUES (:email,:first_name,:last_name,:user_type,:provider,:role_id,:created_at,:updated_at,:password,:is_main_organization_user,:confirmed_at)`,
				map[string]interface{}{
					"email":                     registerUserEmailParam.Email,
					"first_name":                registerUserEmailParam.FirstName,
					"last_name":                 registerUserEmailParam.LastName,
					"user_type":                 registerUserEmailParam.UserType,
					"role_id":                   registerUserEmailParam.RoleId,
					"provider":                  "email",
					"is_main_organization_user": organization == "main_organization",
					"confirmed_at":              time.Now(),
					"password":                  helpers.Functions{}.HashPassword(registerUserEmailParam.Password),
					"created_at":                time.Now(),
					"updated_at":                time.Now(),
				})
		} else {
			_, err = gen.REPO.DB.NamedExec(`INSERT INTO users (email,first_name,last_name,user_type,provider,role_id,created_at,updated_at,password,is_main_organization_user) VALUES (:email,:first_name,:last_name,:user_type,:provider,:role_id,:created_at,:updated_at,:password,:is_main_organization_user)`,
				map[string]interface{}{
					"email":                     registerUserEmailParam.Email,
					"first_name":                registerUserEmailParam.FirstName,
					"last_name":                 registerUserEmailParam.LastName,
					"user_type":                 registerUserEmailParam.UserType,
					"role_id":                   registerUserEmailParam.RoleId,
					"provider":                  "email",
					"is_main_organization_user": organization == "main_organization",
					"password":                  helpers.Functions{}.HashPassword(registerUserEmailParam.Password),
					"created_at":                time.Now(),
					"updated_at":                time.Now(),
				})
		}

	}

	if err != nil {
		fmt.Println(err.Error())

		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error creating user",
		})
		return
	}

	if organization != "main_organization" {
		mail := helpers.Mail{}
		sent := mail.SendMailVerification(registerUserEmailParam.Email, fmt.Sprint(registerUserEmailParam.FirstName, " ", registerUserEmailParam.LastName))
		if sent {
			context.JSON(http.StatusOK, gin.H{
				"error":   false,
				"message": "User registered, email verification link sent to your email",
			})
		} else {
			context.JSON(http.StatusExpectationFailed, gin.H{
				"error":   false,
				"message": "User registered successfully",
			})
		}
	} else {
		context.JSON(http.StatusExpectationFailed, gin.H{
			"error":   false,
			"message": "User registered successfully",
		})
	}

}

func (auth AuthController) ResetPassword(context *gin.Context) {
	var emailResetPasswordParam EmailResetPasswordParam
	err := context.ShouldBindJSON(&emailResetPasswordParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := GetEmailUser(emailResetPasswordParam.Email)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address not found",
		})
		return
	}

	mail := helpers.Mail{}
	sent := mail.SendPasswordResetMail(user.Email.String, fmt.Sprint(user.FirstName.String, " ", user.LastName.String))
	if sent {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Password reset link sent to your email",
		})
	} else {
		context.JSON(http.StatusExpectationFailed, gin.H{
			"error":   true,
			"message": "Error sending password reset email",
		})
	}

}

func (auth AuthController) ResetPasswordApi(context *gin.Context) {
	var emailResetPasswordParam EmailResetPasswordParam
	err := context.ShouldBindJSON(&emailResetPasswordParam)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	user, err := GetEmailUser(emailResetPasswordParam.Email)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address not found",
		})
		return
	}

	mail := helpers.Mail{}
	sent := mail.SendPasswordResetApiMail(user.Email.String, fmt.Sprint(user.FirstName.String, " ", user.LastName.String))
	if sent {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Password reset link sent to your email",
		})
	} else {
		context.JSON(http.StatusExpectationFailed, gin.H{
			"error":   true,
			"message": "Error sending password reset email",
		})
	}

}

func (auth AuthController) EnterNewPassword(context *gin.Context) {
	token := context.DefaultQuery("token", "-")
	valid, err := IsRecoveryTokenValid(token)

	organization := models.TtnmOrganizationModel{}
	fetchOrganizationError := gen.REPO.DB.Get(&organization, "select name, logo_path, website_url from organizations")
	if fetchOrganizationError != nil {
		organization.Name = ""
		organization.LogoPath = ""
		organization.WebsiteUrl = ""
	}

	if err == nil && valid {
		context.HTML(http.StatusOK, "new_password.html", gin.H{
			"token":                    token,
			"organization_name":        organization.Name,
			"organization_logo_path":   organization.LogoPath,
			"organization_website_url": organization.WebsiteUrl,
		})
	}
}

func (auth AuthController) VerifyEmail(context *gin.Context) {
	token := context.DefaultQuery("token", "-")
	valid, err := IsVerificationTokenValid(token)

	organization := models.TtnmOrganizationModel{}
	fetchOrganizationError := gen.REPO.DB.Get(&organization, "select name, logo_path, website_url from ttnm_organization")
	if fetchOrganizationError != nil {
		organization.Name = ""
		organization.LogoPath = ""
		organization.WebsiteUrl = ""
	}

	if err == nil && valid {
		user := models.User{}
		gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select id,email from users where confirmation_token=?"), token)
		_, err = gen.REPO.DB.NamedExec(`UPDATE users set confirmed_at=:confirmed_at where email=:email`,
			map[string]interface{}{
				"confirmed_at": time.Now(),
				"email":        user.Email,
			})
		context.HTML(http.StatusOK, "account_verify_success.html", gin.H{
			"organization_name":        organization.Name,
			"organization_logo_path":   organization.LogoPath,
			"organization_website_url": organization.WebsiteUrl,
		})
	} else {
		context.HTML(http.StatusOK, "account_verify_error.html", gin.H{
			"organization_name":        organization.Name,
			"organization_logo_path":   organization.LogoPath,
			"organization_website_url": organization.WebsiteUrl,
		})
	}
}

func (auth AuthController) SubmitNewPassword(context *gin.Context) {
	err := context.Request.ParseForm()
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error occured",
		})
		return
	}
	//c.Request.PostForm["emails"])
	token := context.PostForm("token")
	password := context.PostForm("password")

	valid, err := IsRecoveryTokenValid(token)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error occured",
		})
		return
	}
	if !valid {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid token/token expired",
		})
		return
	}
	user := models.User{}
	gen.REPO.DB.Get(&user, gen.REPO.DB.Rebind("select id from users where recovery_token=?"), token)

	_, err = gen.REPO.DB.NamedExec(`UPDATE users SET password=:password where id=:id`, map[string]interface{}{
		"password": helpers.Functions{}.HashPassword(password),
		"id":       user.ID,
	})

	if err != nil {
		logger.Log(TAG, fmt.Sprint("Query error:", err), logger.LOG_LEVEL_ERROR)
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Password reset successful",
		})
	}
}

func (auth AuthController) SubmitNewPasswordApi(context *gin.Context) {
	// ResetPasswordApiParams

	var resetPasswordApiParams ResetPasswordApiParams
	err := context.ShouldBindJSON(&resetPasswordApiParams)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetEmailUser(resetPasswordApiParams.Email)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "User with the given email address does not exists",
		})
		return
	}

	if len(resetPasswordApiParams.Password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Password length should be greater than or equal to 6",
		})
		return
	}

	valid, err := IsRecoveryTokenValid(resetPasswordApiParams.Token)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Error occured",
		})
		return
	}
	if !valid {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid token/token expired",
		})
		return
	}

	user_ := models.User{}
	gen.REPO.DB.Get(&user_, gen.REPO.DB.Rebind("select id from users where recovery_token=?"), resetPasswordApiParams.Token)

	_, err = gen.REPO.DB.NamedExec(`UPDATE users SET password=:password where id=:id`, map[string]interface{}{
		"password": helpers.Functions{}.HashPassword(resetPasswordApiParams.Password),
		"id":       user.ID,
	})

	if err != nil {
		logger.Log(TAG, fmt.Sprint("Query error:", err), logger.LOG_LEVEL_ERROR)
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Password reset successful",
		})
	}
}

func (auth AuthController) PasswordResetAndVerifyOTPPhone(context *gin.Context) {
	var resetPasswordPhoneApiParams ResetPasswordPhoneApiParams
	err := context.ShouldBindJSON(&resetPasswordPhoneApiParams)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	user, err := GetUserFromRecoveryTokenWithValidationVerification(resetPasswordPhoneApiParams.OTPCode, resetPasswordPhoneApiParams.CallingCode, resetPasswordPhoneApiParams.Phone)
	if user == nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	_, err = gen.REPO.DB.NamedExec(`UPDATE users SET password=:password where id=:id`, map[string]interface{}{
		"password": helpers.Functions{}.HashPassword(resetPasswordPhoneApiParams.Password),
		"id":       user.ID,
	})
	if err != nil {
		logger.Log(TAG, fmt.Sprint("Query error:", err), logger.LOG_LEVEL_ERROR)
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Pin reset success",
		})
	}
}

func (auth AuthController) PassWordResetSuccess(context *gin.Context) {
	organization := models.TtnmOrganizationModel{}
	err := gen.REPO.DB.Get(&organization, "select name, logo_path, website_url from organizations")
	if err != nil {
		organization.Name = ""
		organization.LogoPath = ""
		organization.WebsiteUrl = ""
	}
	context.HTML(http.StatusOK, "password_reset_success.html", gin.H{
		"organization_name":        organization.Name,
		"organization_logo_path":   organization.LogoPath,
		"organization_website_url": organization.WebsiteUrl,
	})
}
