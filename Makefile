.PHONY: clean test

clean:
	@rm -rf bin/

test: clean
	@./test.sh
