package app

type CliOptions interface {
	// AddFlags adds flags to the specified FlagSet object.
	// AddFlags(fs *pflag.FlagSet)
	Flags() (fss NamedFlagSets)
	Validate() []error
}

// CompleteableOptions abstracts options which can be completed.
type CompleteableOptions interface {
	Complete() error
}

// PrintableOptions abstracts options which can be printed.
type PrintableOptions interface {
	String() string
}
