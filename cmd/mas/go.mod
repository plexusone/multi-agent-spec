module github.com/agentplexus/multi-agent-spec/cmd/mas

go 1.24

require (
	github.com/agentplexus/multi-agent-spec/sdk/go v0.5.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.1
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// For local development - remove before release
replace github.com/agentplexus/multi-agent-spec/sdk/go => ../../sdk/go
