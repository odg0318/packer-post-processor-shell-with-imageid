package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	amazon "github.com/mitchellh/packer/builder/amazon/ebs"
	"github.com/mitchellh/packer/builder/googlecompute"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	Script string `mapstructure:"script"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	if err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
	}, raws...); err != nil {
		return err
	}

	errs := new(packer.MultiError)
	if len(p.config.Script) == 0 {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Error configure: script is required."))
	}

	if _, err := os.Stat(p.config.Script); err != nil {
		errs = packer.MultiErrorAppend(
			errs, fmt.Errorf("Bad script '%s': %s", p.config.Script, err))
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	isRelativePath := string(p.config.Script[0]) != "/"
	if isRelativePath {
		p.config.Script = fmt.Sprintf("./%s", p.config.Script)
	}

	return nil
}

func (p *PostProcessor) ValidBuilder(builderId string) bool {
	return builderId == googlecompute.BuilderId || builderId == amazon.BuilderId
}

func (p *PostProcessor) GetImageId(artifact packer.Artifact) string {
	switch artifact.BuilderId() {
	case googlecompute.BuilderId:
		return artifact.State("ImageName").(string)

	case amazon.BuilderId:
		r, _ := regexp.Compile("ami-[a-z0-9]+")
		amiId := r.FindString(artifact.Id())
		return amiId
	default:
		return ""
	}
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	if p.ValidBuilder(artifact.BuilderId()) == false {
		err := fmt.Errorf("Unknown artifact type: %s.", artifact.BuilderId())
		return nil, false, err
	}
	imageId := p.GetImageId(artifact)

	ui.Message(fmt.Sprintf("Post processing.. image id %s.", imageId))

	cmd := exec.Command(p.config.Script, p.config.PackerBuilderType, imageId)
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		err := fmt.Errorf("Running script failed.")
		return artifact, false, err
	}

	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	if waitStatus.ExitStatus() != 0 {
		err := fmt.Errorf("Running script failed. %d", waitStatus.ExitStatus())
		return artifact, false, err
	}

	ui.Message(fmt.Sprintf("Post processing.. %s done! exit(%d).", p.config.Script, waitStatus.ExitStatus()))

	return artifact, true, nil
}
