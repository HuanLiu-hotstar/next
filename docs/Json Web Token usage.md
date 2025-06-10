# Json Web Token usage



## form 

- Header
- Payload
- Signature



```go
jwt = header.payload.signature
like this xxx.yyyy.zzz

```





### header 

```json
data = {
	"alg":"HS256",
	"typ":"JWT"
}
header = Base64Url(data)
```



### Payload

- contains claims
- Claims are statements about an entity ,typically user data informations
- types of claims : `registered`, `public`,`private`



``` 
Registered claims
{
	"iss" :"issuer",
	"exp":"expiration time",
	"sub":"subject",
	aud:"audience",
	“others”:"others data"  // user defined field
}
```



```
Public claims

```





## How JWT works

- user login will get a token
- client should send token to server with JWT, format maybe like this

``` go
Authorization: Bearer <token>
```

- cross-origin resource sharing will not be an issue 

### server verify the token

- client send request to authorization server , 
- authorization is granted ,then authorization server returns an access token to client 
- client use access token to access API in server