build:
	nex -o src/lex.nn.go -s src/lex.nex
	go build -o gl -i src/main.go src/lex.nn.go src/parse.go
