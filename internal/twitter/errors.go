package twitter

import (
	"fmt"
	spl "snsGoSDK/internal/spl-name-services"
)

var (
	ErrInvalidReverseTwitter spl.SNSError = fmt.Errorf("SNSError: InvalidReverseTwitter")
)
