package algorand

type Capability string

const (
	// Archival nodes
	Archival Capability = "archival"
	// Catchpoints storing nodes
	Catchpoints Capability = "catchpointStoring"
	// Gossip nodes are non permissioned relays
	Gossip Capability = "gossip"
)

const (
	AlgorandKadProtocolID = "/algorand/kad/mainnet/kad/1.0.0"
	AlgorandWsProtocolV22 = "/algorand-ws/2.2.0"
	AlgorandWsProtocolV1  = "/algorand-ws/1.0.0"
)

var advertiseList = []Capability{
	Gossip,
}

var MainNetBootStrapList = []string{
	//host -t txt _dnsaddr.mainnet.algorand.network
	"/dns4/r-jg.algorand-mainnet.network/tcp/4190/p2p/12D3KooWCMWjHGwi8azNw4RcH2FaRZ5oRNF9WL4SQhLBeScy6mii",
	"/dns4/r-g9.algorand-mainnet.network/tcp/4190/p2p/12D3KooWLkYabmTorharKbZsryzwqkkAGCJAsvE5DUwPBLqZjAWF",
	"/dns4/r-gd.algorand-mainnet.network/tcp/4190/p2p/12D3KooWPULa6pcw2hiyLgd72J9HJ49rj6dZnLr6U9e281U2zf9u",
	"/dns4/r-65.algorand-mainnet.network/tcp/4190/p2p/12D3KooWHcrbtg312Rn1L1Z5VytSt6pJc4spMVp9DEPBS18ZYrhw",
	"/dns4/r-8o.algorand-mainnet.network/tcp/4190/p2p/12D3KooWK7NV2QpMu8T8qjhyH1DTjMmHwDAHcX296VYcgy3YBZ1C",
	"/dns4/r-i2.algorand-mainnet.network/tcp/4190/p2p/12D3KooWDXq3KmAwBkTKcj4b6fDYnFVBGUb7ottHkst6mvy7dnN1",
	"/dns4/r-oi.algorand-mainnet.network/tcp/4190/p2p/12D3KooWRavLKs8CwA3ZKnNQJPdjyMYdQzCh7P3U33DArTkQe4XX",
	"/dns4/r-sa.algorand-mainnet.network/tcp/4190/p2p/12D3KooWSdurxomMn38Fk6aqXY1FgWmMnththYxUXX8BNz8N5Kyk",
	"/dns4/r-no.algorand-mainnet.network/tcp/4190/p2p/12D3KooWSAhvCxboKY1abGhUpxzAq6ByC7FbwxwimEgcJjqFziyv",
	"/dns4/r-fj.algorand-mainnet.network/tcp/4190/p2p/12D3KooWA8F3LXxvEiCg1URqFfDU91wW3AtXH1f1WFESp4E8o6sS",
	"/dns4/r-2o.algorand-mainnet.network/tcp/4190/p2p/12D3KooWJaD1NyhzWEdpUqK4zsmN1RhawvPe1hTkKCnJAoCo9dRX",
	"/dns4/r-jy.algorand-mainnet.network/tcp/4190/p2p/12D3KooWRAqPqP9fb4Y3MoSuReXt6DfZeKgmUnb8vW1pKiVz8VjB",
	"/dns4/r-es.algorand-mainnet.network/tcp/4190/p2p/12D3KooWAJRd87K2Zm4k9m49u1ckJgmoGmfAGywS2Xa2mSzxMpui",
	"/dns4/r-iq.algorand-mainnet.network/tcp/4190/p2p/12D3KooWLc4nmyT4AJFcZWpbGW1G4pWmvshZnNDZtuEDe6PaykYj",
	"/dns4/r-cp.algorand-mainnet.network/tcp/4190/p2p/12D3KooWE1sCUiSovmt6evN4idknnVqDfjeUyt8wGwqC4po4fJJN",
	"/dns4/r-53.algorand-mainnet.network/tcp/4190/p2p/12D3KooWQcKxektH8B28MDX7onsP2tKAvpapZquF7s7meAdYmfSE",
	"/dns4/r-d9.algorand-mainnet.network/tcp/4190/p2p/12D3KooWPigdEj6p7NaaMiG2LcN6S5nNg32MPfFY9hZcTXUfgJ6e",
	"/dns4/r-5o.algorand-mainnet.network/tcp/4190/p2p/12D3KooWECRNeewiYkvBBBCdpDCmuqDqKXK2qoFPL1qKsc156jY4",
	"/dns4/r-nh.algorand-mainnet.network/tcp/4190/p2p/12D3KooWRGXhPwFi2BaU3DFpD5AtpfHsgVCfgBpZ5DqVps1tjw2Z",
	"/dns4/r-ab.algorand-mainnet.network/tcp/4190/p2p/12D3KooWGn3qX1NgBkQzHgpmFHMspwaWyG3ioNVwYBYjsrw7YQgD",
	"/dns4/r-s5.algorand-mainnet.network/tcp/4190/p2p/12D3KooWLm9KfdDC9VNmfBdLpTozHbBYhKyoyBL6a9efi4pRW1Cf",
	"/dns4/r-n8.algorand-mainnet.network/tcp/4190/p2p/12D3KooWRJTTmuF1q87nWDpuUUBUGk36cj6oawsxhDTJix8ZiwLD",
	"/dns4/r-fl.algorand-mainnet.network/tcp/4190/p2p/12D3KooWK76LorP7bkPh9VWZNB6RuWqToenPqXMBSgfq4cuUEXrm",
	"/dns4/r-8k.algorand-mainnet.network/tcp/4190/p2p/12D3KooWQ1r5p5YeWTu5sdWBkjqaTu5ewko7Fc1HeetgWXGFtLFo",
	"/dns4/r-xx.algorand-mainnet.network/tcp/4190/p2p/12D3KooWB17vrusb66HREfYwHNxTp8Xcy9sunbBykeDWgp2qRKu4",
	"/dns4/r-kl.algorand-mainnet.network/tcp/4190/p2p/12D3KooWCPKeQ5USNBseZMLQnSD5V5nifUDoMdr2NNcJCn5QAbWt",
	"/dns4/r-pi.algorand-mainnet.network/tcp/4190/p2p/12D3KooWSLn89iUpESaxrXSv4uB8YemGDeKMKShAPUXL5pGA1YyX",
	"/dns4/r-bg.algorand-mainnet.network/tcp/4190/p2p/12D3KooWCQgWskZBa2MNSiN9ZFJUEN1NwfKk4MqqqFKdqr6gNhYe",
	"/dns4/r-bf.algorand-mainnet.network/tcp/4190/p2p/12D3KooWMGhkw7TNGXN9oCqdyv3uUEQwY8Ret6b2mNWoF4yFWHjT",
	"/dns4/r-bk.algorand-mainnet.network/tcp/4190/p2p/12D3KooWRCe8W7zRrdxeFc8c8sJwNUDMPqVEHN9PWT74xaRntt4f",
	"/dns4/r-eh.algorand-mainnet.network/tcp/4190/p2p/12D3KooWPtoCdKuuChZCUPUUEkwy6eLbpVdhhCpRvMtWj4GfAwre",
	"/dns4/r-j5.algorand-mainnet.network/tcp/4190/p2p/12D3KooWCEzCGYcCcimHVYCzipgFYThSBY261T7QMjvkfqBShm97",
	"/dns4/r-kf.algorand-mainnet.network/tcp/4190/p2p/12D3KooWHTBr7xE1xcukR63TV9CkWBby99mzf84A8qzGbWMoSc8a",
	"/dns4/r-9e.algorand-mainnet.network/tcp/4190/p2p/12D3KooWCRn1MnzHzPK1MkDGHRLLWmpDhzW7525nSfWnD46b8zGm",
	"/dns4/r-38.algorand-mainnet.network/tcp/4190/p2p/12D3KooWMBGG9jAAuReNmv2GA8RrD4TPgnTMdh8hehg3EnBSXN7A",
	"/dns4/r-tl.algorand-mainnet.network/tcp/4190/p2p/12D3KooWL77BNRaWrNnBkmfNw4Ja3J5NVjmMHgzEwoZr9Jp1F7ox",
	"/dns4/r-l3.algorand-mainnet.network/tcp/4190/p2p/12D3KooWSMVdcyLSjuiFX3739tP1wMHBhhczCoZ4oZ87k1AxH5Nx",
	"/dns4/r-ej.algorand-mainnet.network/tcp/4190/p2p/12D3KooWKA3okxQeAZzkpidefyHzQf8bmYcKcGY9bKfDTHYawRx9",
	"/dns4/r-r7.algorand-mainnet.network/tcp/4190/p2p/12D3KooWCXXAMvedvh5hqUMZvAAKvCHzgLongQthegHycGp3GqzV",
	"/dns4/r-g5.algorand-mainnet.network/tcp/4190/p2p/12D3KooWQj5dBj9r528tJ1FwuaK9vBUebErnYUb8XGU7qgtvpMPK",
	"/dns4/r-ol.algorand-mainnet.network/tcp/4190/p2p/12D3KooWMzUVfTKJvCGRNHCSuLj3KcmYDBE9ddbTJtfwoH32UQk5",
	"/dns4/r-xw.algorand-mainnet.network/tcp/4190/p2p/12D3KooWMCaUSV4NQzGCADeLVwBsSMQMBBLdVjxCevN5ihbFXSeq",
	"/dns4/r-b0.algorand-mainnet.network/tcp/4190/p2p/12D3KooWEssxxzRAwg8AgABvRoJqrH3ngn27buZsqGSSP7d4m8DF",
	"/dns4/r-rd.algorand-mainnet.network/tcp/4190/p2p/12D3KooWRDzvb2fG8wiQBCTFFkck57yeZpodAgoExgLCrABvgZn2",
	"/dns4/r-aa.algorand-mainnet.network/tcp/4190/p2p/12D3KooWSH9L2dT6GYUg2HKXDzNiY3tjH75ibU4yJsWVWGaQRZ2x",
}
