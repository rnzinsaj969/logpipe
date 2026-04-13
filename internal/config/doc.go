// Package config handles loading and validating logpipe configuration files.
//
// Configuration is expressed as YAML and specifies one or more log sources,
// optional filter criteria (level, service), and output formatting options.
//
// Example configuration:
//
//	sources:
//	  - path: /var/log/app/api.log
//	    service: api
//	  - path: /var/log/app/worker.log
//	    service: worker
//	filters:
//	  level: warn
//	  service: api
//	output:
//	  format: json
//	  destination: /tmp/filtered.log
//
// The Load function returns a validated Config ready for use by the pipeline.
package config
