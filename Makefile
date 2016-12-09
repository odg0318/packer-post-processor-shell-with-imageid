install:
	@go build
	@cp packer-post-processor-shell-with-imageid ~/.packer.d/plugins
