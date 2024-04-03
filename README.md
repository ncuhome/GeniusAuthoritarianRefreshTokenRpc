# Deprecated

rpc will be updated to mutual ssl mode and migrate into client

### Usage

```go
client, err := refreshTokenRpc.NewRpc("v.ncuos.com:443", &refreshTokenRpc.Config{
	AppCode: "xxx",
	AppSecret: "xxx"
})
if err!=nil {
	panic(err)
}

token, err := client.RefreshToken(context.Background(), "refresnTokenHere")
if err!=nil {
	panic(err)
}

fmt.Println(token)
```
