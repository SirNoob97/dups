.PHONY: clean test

clean:
	@rm -rf bin/ testdata/

test: clean
	@./test.sh
