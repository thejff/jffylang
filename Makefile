all: ast build

cleanbuild: clean all

build:
	# Build JFFY lang
	go build -o bin/jffy main.go

ast:
	# Build AST Generator
	go build -o bin/generateast tool/generateast/main.go

	# Generate ASTs
	./bin/generateast ./jffy

clean: 
	rm -rf ./bin/*

