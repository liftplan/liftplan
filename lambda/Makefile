project_name=liftplan
project_dir=function
builds_dir=builds
terraform_dir=terraform
TF_BACKEND_FILE?=$(HOME)/liftplan-lambda-backend.hcl

init:
	@cd $(terraform_dir) && \
	terraform init -backend-config=$(TF_BACKEND_FILE)

clean:
	@rm -rf $(project_dir)/$(builds_dir)/*

build: clean
	@cd $(project_dir) && \
	GOOS=linux go build -o $(builds_dir)/$(project_name) main.go

zip: build
	@cd $(project_dir)/$(builds_dir) && \
	zip function.zip $(project_name)

