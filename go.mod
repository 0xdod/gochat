module github.com/0xdod/gochat

go 1.16

// +heroku goVersion go1.16
// +heroku install ./cmd/...


require (
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible // indirect
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/jackc/pgconn v1.8.0
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/satori/uuid v1.2.0 // indirect
	github.com/stretchr/objx v0.3.0
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)
