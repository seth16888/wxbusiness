package jwt

import (
	"errors"

	"time"

	goJWT "github.com/golang-jwt/jwt/v5"
	"github.com/seth16888/wxbusiness/pkg/helpers"
	"go.uber.org/zap"
)

var (
	ErrTokenExpired           = errors.New("令牌过期")
	ErrTokenNotValidYet       = errors.New("令牌未生效")
	ErrTokenMalformed         = errors.New("令牌格式错误")
	ErrTokenSignatureInvalid  = errors.New("令牌签名无效")
	ErrTokenInvalid           = errors.New("请求令牌无效")
	ErrTokenExpiredMaxRefresh = errors.New("令牌已过最大刷新时间")
)

// JWTService JWT服务
type JWTService struct {
	SignKey        string
	Issuer         string
	ExpireTime     time.Duration
	MaxRefreshTime time.Duration
	TimeLocation   *time.Location

	logger *zap.Logger
}

func (j *JWTService) GetExpires() time.Duration {
	return j.ExpireTime
}

type JWTCustomClaims struct {
	// userID
	UserID int64 `json:"userId"`
	// authorities 权限列表
	Authorities  []string `json:"authorities"`
	DepartmentID int64    `json:"deptId"`
	DataScope    int      `json:"dataScope"`

	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	goJWT.RegisteredClaims
}

// NewJWTService 创建JWT服务
func NewJWTService(signKey string, issuer string, expireTime time.Duration,
	maxRefreshTime time.Duration, timeLocation *time.Location, logger *zap.Logger) *JWTService {
	// 不允许 signKey 为空
	if signKey == "" {
		panic("JWT 密钥不能为空")
	}

	return &JWTService{
		SignKey:        signKey,
		Issuer:         issuer,
		ExpireTime:     expireTime,
		MaxRefreshTime: maxRefreshTime,
		TimeLocation:   timeLocation,
		logger:         logger,
	}
}

// IssueToken 生成 Token
func (j *JWTService) IssueToken(claims *JWTCustomClaims) (tokenString string, exp int64, err error) {
	if claims == nil {
		return "", 0, errors.New("claims is nil")
	}

	tokenString, err = j.createToken(claims)
	if err != nil {
		j.logger.Error("JWT", zap.Any("IssueToken", err.Error()))
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt.Unix(), nil
}

// createToken 创建 Token
func (j *JWTService) createToken(claims *JWTCustomClaims) (tokenString string, err error) {
	token := goJWT.NewWithClaims(goJWT.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SignKey))
}

// ParseToken 解析 Token
func (j *JWTService) ParseToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := j.parseTokenString(tokenString)
	if err != nil {
		switch {
		case errors.Is(err, goJWT.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, goJWT.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		case errors.Is(err, goJWT.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, goJWT.ErrTokenSignatureInvalid):
			return nil, ErrTokenSignatureInvalid
		default:
			return nil, err
		}
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新 Token
//
//	刷新令牌时，需要验证旧令牌，同时生成新的令牌
//	旧令牌通过验证后，生成新的令牌
//	旧令牌不通过验证，返回 ErrTokenInvalid
//	最大刷新时间：从 Token 的签名时间 + jwt.max_refresh_time(分钟)
//	签名时间：第一次签发时间，后续刷新不改变，只改变过期时间
//	刷新Token逻辑：
//	 1. 旧令牌未过期，且未达到最大刷新时间，返回新令牌
//	 2. 旧令牌未过期，但已达到最大刷新时间，返回 ErrTokenExpiredMaxRefresh
//	 3. 旧令牌过期，返回 ErrTokenExpired
func (j *JWTService) RefreshToken(tokenString string) (token string, err error) {
	jwt, err := j.parseTokenString(tokenString)
	if err != nil {
		return "", err
	}

	// claims
	claims := jwt.Claims.(*JWTCustomClaims)

	// 刷新最大时间
	maxRefreshTime := claims.IssuedAt.Time.Add(j.MaxRefreshTime)

	// 过期返回错误
	if helpers.TimenowInTimezone(j.TimeLocation).After(maxRefreshTime) {
		return "", ErrTokenExpiredMaxRefresh
	}

	// 可以刷新，重新签发
	// 过期时间
	t := helpers.TimenowInTimezone(j.TimeLocation).Add(j.ExpireTime)
	claims.RegisteredClaims.ExpiresAt = &goJWT.NumericDate{Time: t}
	return j.createToken(claims)
}

// parseTokenString 解析 Token 字符串
func (j *JWTService) parseTokenString(tokenString string) (*goJWT.Token, error) {
	return goJWT.ParseWithClaims(
		tokenString,
		&JWTCustomClaims{},
		func(token *goJWT.Token) (interface{}, error) {
			return []byte(j.SignKey), nil
		})
}
