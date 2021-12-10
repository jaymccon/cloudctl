package aws

import (
	"context"
	"fmt"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const (
	appName    = "cloudctl"
	version    = "0.0.0-alpha.1"
	appComment = "+https://github.com/jaymccon/cloudctl"
)

var agent = fmt.Sprintf("%s/%s (%s)", appName, version, appComment)

var customerUAMiddleware = middleware.BuildMiddlewareFunc("CustomerUserAgent", func(
	ctx context.Context, input middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	request, ok := input.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", input.Request)
	}

	const userAgentKey = "User-Agent"

	value := request.Header.Get(userAgentKey)

	if len(value) > 0 {
		value = agent + " " + value
	} else {
		value = agent
	}

	request.Header.Set(userAgentKey, value)

	return next.HandleBuild(ctx, input)
})

func attachCustomMiddleware(stack *middleware.Stack) error {
	return stack.Build.Add(customerUAMiddleware, middleware.After)
}
