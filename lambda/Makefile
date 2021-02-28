project_name=liftplan
project_dir=function
builds_dir=builds
terraform_dir=terraform
zip_name=function.zip
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
	zip $(zip_name) $(project_name)

push-to-s3: zip
	@cd $(project_dir)/$(builds_dir) && \
	aws s3 cp $(zip_name) s3://$(LIFTPLAN_LAMBDA_S3_BUCKET)/$(project_name)/$(zip_name)

update: push-to-s3
	@aws lambda update-function-code \
    --function-name  $(project_name) \
    --s3-bucket $(LIFTPLAN_LAMBDA_S3_BUCKET) \
	--s3-key $(project_name)/$(zip_name)