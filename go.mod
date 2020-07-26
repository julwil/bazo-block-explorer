module github.com/julwil/bazo-block-explorer

go 1.14

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/julwil/bazo-miner v0.0.0-20200707114909-5dc078567dad
	github.com/lib/pq v1.7.1
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
)

replace github.com/julwil/bazo-miner => ../bazo-miner // Packages from bazo-miner are resolved locally, rather than with the specified version.