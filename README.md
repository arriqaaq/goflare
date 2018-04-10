# goflare

A simple command line tool to query cloudflare dns


## Usage

#Resolve

	./goflare -name="duckduckgo.com" --action="resolve" --qtype=1
	./goflare -name="duckduckgo.com" --action="resolve" --qtype=AAAA


#Query

	./goflare -name="duckduckgo.com" --action="resolve" --qtype=1
	./goflare -name="duckduckgo.com" --action="resolve" --qtype=AAAA