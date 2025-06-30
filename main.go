package main

import (
	"github.com/chuhaoyuu/aws-oidc-sts/cmd"
	_ "github.com/chuhaoyuu/aws-oidc-sts/cmd/aws"
)

func main() {
	cmd.Execute()
}
