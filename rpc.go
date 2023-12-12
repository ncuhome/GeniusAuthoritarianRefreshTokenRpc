package refreshTokenRpc

import (
	"context"
	"crypto/sha256"
	"fmt"
	refreshTokenProto "github.com/ncuhome/GeniusAuthoritarianProtos/refreshToken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"unsafe"
)

func NewRpc(addr string, c *Config) (*Rpc, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}
	client := refreshTokenProto.NewRefreshTokenClient(conn)
	return &Rpc{
		conn: client,
		conf: c,
	}, nil
}

type Rpc struct {
	conn refreshTokenProto.RefreshTokenClient
	conf *Config
}

func (r Rpc) FillSignature(req *refreshTokenProto.TokenRequest) *refreshTokenProto.TokenRequest {
	req.AppCode = r.conf.AppCode
	signStr := fmt.Sprintf("%s:%s:%s", req.AppCode, r.conf.AppSecret, req.Token)

	h := sha256.New()
	h.Write(unsafe.Slice(unsafe.StringData(signStr), len(signStr)))
	req.Signature = fmt.Sprintf("%x", h.Sum(nil))

	return req
}

func (r Rpc) RefreshToken(ctx context.Context, refreshToken string) (token *refreshTokenProto.AccessToken, err error) {
	return r.conn.RefreshToken(ctx, r.FillSignature(&refreshTokenProto.TokenRequest{
		Token: refreshToken,
	}))
}

func (r Rpc) DestroyToken(ctx context.Context, refreshToken string) error {
	_, err := r.conn.DestroyRefreshToken(ctx, r.FillSignature(&refreshTokenProto.TokenRequest{
		Token: refreshToken,
	}))
	return err
}

func (r Rpc) VerifyAccessToken(ctx context.Context, accessToken string) (*refreshTokenProto.AccessTokenInfo, error) {
	return r.conn.VerifyAccessToken(ctx, r.FillSignature(&refreshTokenProto.TokenRequest{
		Token: accessToken,
	}))
}

func (r Rpc) GetUserInfo(ctx context.Context, accessToken string) (*refreshTokenProto.UserInfo, error) {
	return r.conn.GetUserInfo(ctx, r.FillSignature(&refreshTokenProto.TokenRequest{
		Token: accessToken,
	}))
}
