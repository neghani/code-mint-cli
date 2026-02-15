package cobra

import "testing"

func TestCommandParsesFlagsAfterPositionalArgs(t *testing.T) {
	var tool string
	var dryRun bool
	var gotArgs []string

	cmd := &Command{
		Use: "add",
		RunE: func(_ *Command, args []string) error {
			gotArgs = args
			return nil
		},
	}
	cmd.Flags().StringVar(&tool, "tool", "", "")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "")

	if err := cmd.execute([]string{"@rule/safe-api-route-pattern", "--tool", "cursor", "--dry-run"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}

	if len(gotArgs) != 1 || gotArgs[0] != "@rule/safe-api-route-pattern" {
		t.Fatalf("unexpected positional args: %#v", gotArgs)
	}
	if tool != "cursor" {
		t.Fatalf("tool not parsed, got %q", tool)
	}
	if !dryRun {
		t.Fatalf("dry-run not parsed")
	}
}

func TestNormalizeInterspersedArgsRespectsDoubleDash(t *testing.T) {
	cmd := &Command{
		Use: "cmd",
		RunE: func(_ *Command, _ []string) error {
			return nil
		},
	}
	var debug bool
	cmd.Flags().BoolVar(&debug, "debug", false, "")

	if err := cmd.execute([]string{"value", "--", "--debug"}); err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if debug {
		t.Fatalf("flag after -- should not be parsed")
	}
}
