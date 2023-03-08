package options

import (
	"encoding/json"
	"golang-standards-project-example/internal/pkg/options"
	"golang-standards-project-example/pkg/app"
)

type Options struct {
	HttpServingOptions *options.HttpServingOptions `json:"http"   mapstructure:"http"`
}

func NewOptions() *Options {
	return &Options{
		HttpServingOptions: options.NewHttpServingOptions(),
	}
}

func (o *Options) Flags() (fss app.NamedFlagSets) {
	o.HttpServingOptions.AddFlags(fss.FlagSet("http"))
	return
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.HttpServingOptions.Validate()...)

	return errs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
