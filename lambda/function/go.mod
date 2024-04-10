require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/liftplan/liftplan v0.0.2
)

replace github.com/liftplan/liftplan => ../../

module github.com/liftplan/liftplan/lambda/runtime

go 1.22
