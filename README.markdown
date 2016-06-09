# Go Project Template

Using bootstrap, vue.js, gorp (postgres), and jsontokens for auth.

# Installation

Install go 1.6.2 or later

## Dependencies

Install following depdendences:

go get github.com/gorilla/mux
go get github.com/gorilla/sessions
go get github.com/lib/pq
go get bitbucket.org/liamstask/goose/cmd/goose
go get golang.org/x/crypto/bcrypt
go get github.com/goods/httpbuf
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
go get github.com/sclevine/agouti
go get gopkg.in/gorp.v1
go get github.com/dgrijalva/jwt-go
go get github.com/kardianos/osext

To update dependencies:

go get -u github.com/gorilla/mux
go get -u github.com/gorilla/sessions
go get -u github.com/lib/pq
go get -u bitbucket.org/liamstask/goose/cmd/goose
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/goods/httpbuf
go get -u github.com/onsi/ginkgo/ginkgo
go get -u github.com/onsi/gomega
go get -u github.com/sclevine/agouti
go get -u gopkg.in/gorp.v1
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/kardianos/osext

The following are needed to run tests:

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
go get github.com/sclevine/agouti

To update:

go get -u github.com/onsi/ginkgo/ginkgo
go get -u github.com/onsi/gomega
go get -u github.com/sclevine/agouti

== Setting up secure keys ==

Generate the secure keys as follows.

1.  Install pwgen

sudo apt-get install pwgen

2.  Generate the keys for the cookie store

  i.   Generate an authentication key
         pwgen -s 64 -1 -n
  ii.  Add the result of the above command to ~/.bashrc (or platform equivalent) (replace example passkey with output of 2.i.)
         export GO_TEMPLATE_COOKIE_STORE_AUTHENTICATION_KEY=r4njzugGZk1esyzXT9CtqW5rYJaNex7J0EQtjfgM6wlwXr9kNio4hp9XjZiYeamr
  iii. Generate an encryption key
         pwgen -s 32 -1 -n
  iv.  Add the result of the above command to ~/.bashrc (or platform equivalent) (replace example passkey with output of 2.i.)
         export GO_TEMPLATE_COOKIE_STORE_ENCRYPTION_KEY=r4njzugGZk1esyzXT9CtqW5rYJaNex7J0EQtjfgM6wlwXr9kNio4hp9XjZiYeamr

3.  Generate JWT a JWT signing key

  i.   generate the key
         pwgen -s 32 -1 -n
  ii.  add the result of the above command ~/.bashrc (or platform equivalent) (replace example passkey with output of 2.i.)
         export GO_TEMPLATE_JWT_KEY=sQyYUSy52odvX5Qdxt5SV1RsjuXEw7nm        

== Setting up the database ==

1.  Launch the postgres command line
    psql -U postgres -W
2.  Create the database
    CREATE DATABASE gotemplate;
  i. Repeat for developement and test databases as needed.
3.  Create the user:
    CREATE USER gotemplate WITH PASSWORD 'password';
4.  Grant the new user access to the database:
    GRANT ALL PRIVILEGES ON DATABASE gotemplate TO gotemplate;
  i.  Repeat for development and test databases as needed.
5.  Exit postgres
    \q
6.  Run the goose migrations
    goose up -env production
  i.  Repeat for development and test databases as needed

