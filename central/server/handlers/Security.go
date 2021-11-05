package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/newm4n/mihp/central/model"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	Issuer        = "mihp"
	SignKey       = "bsa2fe6h5j3a3s5k8e8fj5h3v3oa33i3u7ck0be"
	AccessKeyAge  = 10 * time.Minute
	RefreshKeyAge = 24 * 30 * 12 * time.Hour

	userRepo model.UserRepository
)

func ValidateAuthorizationToken(ctx *fasthttp.RequestCtx) bool {
	auth := ctx.Request.Header.Peek("Authorization")
	if auth == nil {
		return false
	}
	sauth := string(auth)
	if len(auth) <= 7 || !strings.HasPrefix(strings.ToUpper(sauth), "BEARER ") {
		return false
	}
	//token := sauth[7:]
	return true
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, "missing body", fasthttp.StatusBadRequest)
		return
	} else {
		contentType := "ctx.Request.Header.ContentType()"
		if contentType != "application/json" {
			ErrorResponse(w, "Content-Type not json", fasthttp.StatusBadRequest)
			return
		}
		reqBody := &LoginRequest{}
		err := json.Unmarshal(bodyBytes, reqBody)
		if err != nil {
			ErrorResponse(w, "Error while parsing json body. got "+err.Error(), fasthttp.StatusBadRequest)
			return
		}
		//roles, err := userRepo.GetUserRole(reqBody.Email, reqBody.Password)
		//if err != nil {
		//	logrus.Errorf("got error while retrieving user role. got %s", err.Error())
		//	ErrorResponse(w,"unauthorized", fasthttp.StatusUnauthorized)
		//	return
		//}

	}
}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type RefreshResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"access"`
}

// ReadJWTStringToken takes a token string , keys, signMethod and returns its content.
func ReadJWTStringToken(validate bool, signKey, signMethod, tokenString string) (*JWTSpec, error) {
	if signKey == "th15mustb3CH@ngedINprodUCT10N" {
		logrus.Warnf("Using default CryptKey for JWT Token, This key is visible from the source tree and to be used in development only. YOU MUST CHANGE THIS IN PRODUCTION or TO REMOVE THIS LOG FROM APPEARING")
	}

	jwt, err := jws.ParseJWT([]byte(tokenString))
	if err != nil {
		return nil, fmt.Errorf("malformed jwt token")
	}

	if validate {
		var sMethod crypto.SigningMethod

		switch strings.ToUpper(signMethod) {
		case "HS256":
			sMethod = crypto.SigningMethodHS256
		case "HS384":
			sMethod = crypto.SigningMethodHS384
		case "HS512":
			sMethod = crypto.SigningMethodHS512
		default:
			sMethod = crypto.SigningMethodHS256
		}

		if err := jwt.Validate([]byte(signKey), sMethod); err != nil {
			return nil, fmt.Errorf("invalid jwt token - %s", err.Error())
		}
	}
	claims := jwt.Claims()
	additional := make(map[string]interface{})
	for k, v := range claims {
		kup := strings.ToUpper(k)
		if kup != "ISS" && kup != "AUD" && kup != "SUB" && kup != "IAT" && kup != "EXP" && kup != "NBF" {
			additional[k] = v
		}
	}

	issuer, _ := claims.Issuer()
	subject, _ := claims.Subject()
	audience, _ := claims.Audience()
	expire, _ := claims.Expiration()
	notBefore, _ := claims.NotBefore()
	issuedAt, _ := claims.IssuedAt()

	spec := &JWTSpec{
		SignKey:    signKey,
		SignMethod: signMethod,
		Issuer:     issuer,
		Subject:    subject,
		Audiences:  audience,
		IssuedAt:   issuedAt,
		NotBefore:  notBefore,
		ExpireAt:   expire,
		Additional: additional,
	}

	return spec, nil
}

type JWTSpec struct {
	SignKey    string
	SignMethod string
	Issuer     string
	Subject    string
	Audiences  []string
	IssuedAt   time.Time
	NotBefore  time.Time
	ExpireAt   time.Time
	Additional map[string]interface{}
}

// CreateJWTStringToken create JWT String token based on arguments
func CreateJWTStringToken(spec *JWTSpec) (string, error) {
	if spec.SignKey == "th15mustb3CH@ngedINprodUCT10N" {
		logrus.Warnf("Using default CryptKey for JWT Token, This key is visible from the source tree and to be used in development only. YOU MUST CHANGE THIS IN PRODUCTION or TO REMOVE THIS LOG FROM APPEARING")
	}

	claims := jws.Claims{}
	claims.SetIssuer(spec.Issuer)
	claims.SetSubject(spec.Subject)
	claims.SetAudience(spec.Audiences...)
	claims.SetIssuedAt(spec.IssuedAt)
	claims.SetNotBefore(spec.NotBefore)
	claims.SetExpiration(spec.ExpireAt)

	for k, v := range spec.Additional {
		claims[k] = v
	}

	var signM crypto.SigningMethod

	switch strings.ToUpper(spec.SignMethod) {
	case "HS256":
		signM = crypto.SigningMethodHS256
	case "HS384":
		signM = crypto.SigningMethodHS384
	case "HS512":
		signM = crypto.SigningMethodHS512
	default:
		signM = crypto.SigningMethodHS256
	}

	jwtBytes := jws.NewJWT(claims, signM)

	tokenByte, err := jwtBytes.Serialize([]byte(spec.SignKey))
	if err != nil {
		panic(err)
	}
	return string(tokenByte), nil
}
