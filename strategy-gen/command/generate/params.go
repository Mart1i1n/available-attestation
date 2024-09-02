package generate

import (
	"errors"
	"github.com/tsinghua-cel/strategy-gen/command/generate/config"
)

const (
	configFlag            = "config"
	validatorCountFlag    = "validator-count"
	startSlotFlag         = "start-slot"
	endSlotFlag           = "end-slot"
	enableAttFlag         = "enable-att-points"
	enableBlockFlag       = "enable-block-points"
	enableAttActionFlag   = "enable-att-actions"
	enableBlockActionFlag = "enable-block-actions"
)

var (
	errInvalidNATAddress = errors.New("could not parse NAT IP address")
)

type toolParam struct {
	configPath   string
	outputFile   string
	generateMode int
	rawConfig    *config.Config
}

var (
	params = &toolParam{
		rawConfig: config.DefaultConfig(),
	}
)

func (p *toolParam) initConfigFromFile() error {
	var parseErr error

	if p.rawConfig, parseErr = config.ReadConfigFile(p.configPath); parseErr != nil {
		return parseErr
	}

	return nil
}
