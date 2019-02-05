# gcp-iap-auth

This project is a copy for tests of "github.com/imkira/gcp-iap-auth". Thanks !

`gcp-iap-auth` is a simple server implementation and package in
[Go](http://golang.org) for helping you secure your web apps running on GCP
behind a
[Google Cloud Platform's IAP (Identity-Aware Proxy)](https://cloud.google.com/iap/docs/) by validating IAP signed headers in the requests.

# Why

Validating signed headers helps you protect your app from the following kinds of risks:

- IAP is accidentally disabled;
- Misconfigured firewalls;
- Access from within the project.

## How to use it as a package

```go
go get github.com/slavayssiere/gcp-iap-auth
```

## Sample code

```go
package main

import (
	"github.com/slavayssiere/gcp-iap-auth"
)

var cfg *jwt.Config

// In here we initialize the configuration for our app.
// It doesn't need to be in "init".
func init() {
	reg, err := regexp.Compile(`/projects/YOURPROJECTID/global/backendServices/BACKENDID$`)
	if err != nil {
		log.Fatal(err)
	}

	publicKeys, err := jwt.FetchPublicKeys()
	if err != nil {
		log.Fatal(err)
	}
	cfg = &jwt.Config{
		PublicKeys:     publicKeys,
		MatchAudiences: reg,
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
}

// LoggerMiddleware add logger and metrics
func LoggerMiddleware(inner http.HandlerFunc, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Compare(name, "root") != 0 {

			start := time.Now()

			log.Println(name)
			//tokenString := (*r).Header["X-Goog-Iap-Jwt-Assertion"][0]
			claims, err := jwt.RequestClaims(r, cfg)
			log.Println(err)
			log.Println(claims)
			log.Println(claims.Email)
			log.Println(claims.Subject)

			// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 	// Don't forget to validate the alg is what you expect:
			// 	if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			// 	}

			// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			// 	return []byte("BTX6xf6J7nBXFCSHEQJ3AChg"), nil
			// })

			// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 	log.Println(claims)
			// } else {
			// 	log.Println(err)
			// }

			inner.ServeHTTP(w, r)

			time := time.Since(start)
			log.Printf(
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				name,
				time,
			)
		} else {
			inner.ServeHTTP(w, r)
		}
	})
}

```
