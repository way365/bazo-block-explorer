module github.com/way365/bazo-block-explorer

go 1.20

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.7.1
	github.com/way365/bazo-miner v0.0.0-20200707114909-5dc078567dad
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
)

require (
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	github.com/willf/bloom v2.0.3+incompatible // indirect
	golang.org/x/sys v0.0.0-20190412213103-97732733099d // indirect
)

replace github.com/way365/bazo-miner => ../bazo-miner // Packages from bazo-miner are resolved locally, rather than with the specified version.
