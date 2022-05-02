package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli"

	"github.com/smallstep/cli/command/ca/policy/policycontext"
	"github.com/smallstep/cli/flags"
	"github.com/smallstep/cli/internal/command"
	"github.com/smallstep/cli/utils/cautils"
)

// CommonNamesCommand returns the common names policy subcommand.
func CommonNamesCommand(ctx context.Context) cli.Command {
	return cli.Command{
		Name:  "cn",
		Usage: "...",
		UsageText: `**cn** <domain> [**--remove**]
[**--provisioner**=<name>] [**--key-id**=<key-id>] [**--reference**=<reference>]
[**--admin-cert**=<file>] [**--admin-key**=<file>]
[**--admin-provisioner**=<string>] [**--admin-subject**=<string>]
[**--password-file**=<file>] [**--ca-url**=<uri>] [**--root**=<file>]
[**--context**=<name>]`,
		Description: `**dns** command group provides facilities for ...`,
		Action: command.InjectContext(
			ctx,
			commonNamesAction,
		),
		Flags: []cli.Flag{
			provisionerFilterFlag,
			flags.KeyID,
			flags.Reference,
			flags.AdminCert,
			flags.AdminKey,
			flags.AdminProvisioner,
			flags.AdminSubject,
			flags.PasswordFile,
			flags.CaURL,
			flags.Root,
			flags.Context,
			cli.BoolFlag{
				Name:  "remove",
				Usage: `removes the provided DNS names from the policy instead of adding them`,
			},
		},
	}
}

func commonNamesAction(ctx context.Context) (err error) {

	clictx := command.CLIContextFromContext(ctx)

	args := clictx.Args()
	if len(args) == 0 {
		return errors.New("please provide at least one name")
	}

	client, err := cautils.NewAdminClient(clictx)
	if err != nil {
		return fmt.Errorf("error creating admin client: %w", err)
	}

	policy, err := retrieveAndInitializePolicy(ctx, client)
	if err != nil {
		return fmt.Errorf("error retrieving policy: %w", err)
	}

	var commonNames []string

	switch {
	case policycontext.HasSSHHostPolicy(ctx):
		return errors.New("SSH host policy does not support Common Names")
	case policycontext.HasSSHUserPolicy(ctx):
		return errors.New("SSH user policy does not support Common Names")
	case policycontext.HasX509Policy(ctx):
		switch {
		case policycontext.HasAllow(ctx):
			commonNames = policy.X509.Allow.CommonNames
		case policycontext.HasDeny(ctx):
			commonNames = policy.X509.Deny.CommonNames
		default:
			panic(errors.New("no allow nor deny context set"))
		}
	default:
		panic("no SSH nor X.509 context set")
	}

	if clictx.Bool("remove") {
		for _, domain := range args {
			commonNames = remove(domain, commonNames)
		}
	} else {
		commonNames = append(commonNames, args...)
	}

	switch {
	case policycontext.HasSSHHostPolicy(ctx):
		return errors.New("SSH host policy does not support Common Names")
	case policycontext.HasSSHUserPolicy(ctx):
		return errors.New("SSH user policy does not support Common Names")
	case policycontext.HasX509Policy(ctx):
		switch {
		case policycontext.HasAllow(ctx):
			policy.X509.Allow.CommonNames = commonNames
		case policycontext.HasDeny(ctx):
			policy.X509.Deny.CommonNames = commonNames
		default:
			panic(errors.New("no allow nor deny context set"))
		}
	default:
		panic("no SSH nor X.509 context set")
	}

	updatedPolicy, err := updatePolicy(ctx, client, policy)
	if err != nil {
		return fmt.Errorf("error updating policy: %w", err)
	}

	prettyPrint(updatedPolicy)

	return nil
}
