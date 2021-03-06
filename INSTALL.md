# Installation
This file describes how to setup an AAS with the example service.
It also contains some bits about how to login/register, etc.

# Dependencies
## Golang
Visit their [install](http://golang.org/doc/install) page.

Tested with:

	(earth)% go version
	go version go1.2.1 linux/386

## PostgreSQL
Have a look at the INSTALL file at the root of the source
archive. It contains clean informations on how to build and launch
a PostgreSQL server.

Alternatively, there's the wiki [install](https://wiki.postgresql.org/wiki/Detailed_installation_guides) page.

Tested with:

	(earth)% postgres --version
	postgres (PostgreSQL) 9.3.4

Create a new 'auth' database and add an 'auth' [role](http://www.postgresql.org/docs/9.4/static/user-manag.html):

	(earth)% createdb auth
	(earth)% psql auth
	psql (9.3.4)
	Type "help" for help.
	
	auth=# CREATE ROLE auth PASSWORD 'auth' LOGIN;
	CREATE ROLE
	auth=# \q

## Go packages
### PostgreSQL driver
Providing a driver to communicate with the PostgreSQL database.

Available at [github.com/lib/pq](https://github.com/lib/pq)

### gorilla/securecookie
Encode, encrypts and manage cookies. Part of the
[Gorilla](http://www.gorillatoolkit.org/) web toolkit.

Available at [github.com/gorilla/securecookie](https://github.com/gorilla/securecookie)

### dchest/captcha
Simple captcha management.

Available at [github.com/dchest/captcha](https://github.com/dchest/captcha/)

# Building
Fetch sources:

	(earth)% git clone https://github.com/heaumer/auth
	...
	(earth)% cd auth

Or

	(earth)% go get github.com/heaumer/auth
	...
	(earth)% cd $GOPATH/src/github.com/heaumer/auth

Ensure dependencies are installed:

	(earth)% go get -v ./...
	github.com/dchest/captcha (download)
	github.com/gorilla/securecookie (download)
	github.com/lib/pq (download)
	github.com/dchest/captcha
	github.com/gorilla/securecookie
	github.com/lib/pq/oid
	github.com/lib/pq
	_/home/mb/src/newsome/auth/example
	_/home/mb/src/newsome/auth

Compile:

	(earth)% go build . && (cd example/; go build .) && (cd storexample/; go build .)

# Configuration
## Requirements
Successful build. Also, grab an email account.

## Generating x509 certificate
Use the genkey.sh script : it will generate a x509 cert/key
pair and install it for auth, example and storexample. It
will also copy the certificate to example/conf/auth-cert.conf
and to storexample/auth-cert.conf

See README.md:/### example\/main.go for a description of
example's configuration.

See README.md:/### storexample\/main.go for a description
of storexample's command line options.

## AAS
### Configuration
Have a look at config.json and read the README section
describing it (README.md:/^## Configuration).
You should mainly be interested in the SMTPServer and
associated fields.

We will assume the default options in the following,
so be carefull.

### Launching and connected as admin
Once you edited config.json to fill your needs, you
may launch the AAS issuing a:

	(earth)% ./auth &
	2014/05/18 14:01:41 Launching on https://localhost:8080

Browse to the indicated URL, try to login as 'admin' : an
email with a token will be sent from AuthEmail to AdminEmail
(config.json fields).

Grab the token and use it to login.

## Example
### Service registration
Now that the AAS is launched, let's register a new service to
the AAS, issuing a request to /api/discover:

	(earth)% curl -k --data 'name=example&url=http://example.awesom.eu/&address=127.0.0.1&email=admin@whatev.er' https://localhost:8080/api/discover

The previous requests should return either

* ko : something wrong happened, or server is in Disable mode
* ok : server is in Manual mode. Go to the admin panel and grab the service's key
* somehexstring : Automatic mode, here's your service key

Alternatively, you may add the service from the admin panel.

### Configuration
Edit example/conf/auth.conf according to your key, the url
of your AAS. Unless you changed something, the cert file should
be ok; you may check it does exists in example/conf/auth-cert.pem.

### Launching
Issue a:

	(earth)% cd example/; ./example/
	2014/05/18 14:16:09 Launching on https://localhost:8082

Browse to the indicated URL. Try to login as 'admin'.
Then, go to your sessions page on the AAS, grab the key
for the service, and login with it on the service's page.

You should see your id, nickname an email, and some user data
if the storexample has been started, and if the user has uploaded
some data there.

You may then logout.

## Store Example
### Service registration
In a similar way,

	(earth)% curl -k --data 'name=storexample&url=http://store.awesom.eu/&address=127.0.0.1&email=admin@whatev.er' https://localhost:8080/api/discover

### Configuration
Check the command lines options, mainly `-key`.

### Launching
Start it with a:

	(earth)% cd storexample/; ./storexample

You may now use storexample's API through the previous service.
