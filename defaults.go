package apex

func defaults(opts MetricsOpts) MetricsOpts {
	// opts.Namespace default is empty
	// opts.Subsystem default is empty
	// opts.MustRegister default is false
	if opts.Path == "" {
		opts.Path = "/metrics"
	}

	if opts.Port == 0 {
		opts.Port = 9000
	}

	if opts.Separator == 0 {
		opts.Separator = '_'
	}

	return opts
}
