TARG=favico

$(TARG): *.go
	@go build -o $(TARG) *.go

test: *.go
	@go test -v *.go

bench: *.go
	@go test -bench=".*"

clean:
	@rm -fr $(TARG)
