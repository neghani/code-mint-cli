package cobra

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Use               string
	Short             string
	Run               func(cmd *Command, args []string)
	RunE              func(cmd *Command, args []string) error
	PersistentPreRunE func(cmd *Command, args []string) error

	parent          *Command
	subcommands     []*Command
	flags           *FlagSet
	persistentFlags *FlagSet
	ctx             context.Context
}

func (c *Command) AddCommand(cmds ...*Command) {
	for _, cmd := range cmds {
		cmd.parent = c
		c.subcommands = append(c.subcommands, cmd)
	}
}

func (c *Command) Execute() error {
	if c.ctx == nil {
		c.ctx = context.Background()
	}
	args := os.Args[1:]
	return c.execute(args)
}

func (c *Command) execute(args []string) error {
	if wantsHelp(args) {
		fmt.Println(c.help())
		return nil
	}
	if len(c.subcommands) > 0 {
		for idx, arg := range args {
			if strings.HasPrefix(arg, "-") {
				continue
			}
			for _, sub := range c.subcommands {
				use := strings.Fields(sub.Use)
				if len(use) > 0 && arg == use[0] {
					if c.persistentFlags != nil && c.persistentFlags.hasAny() {
						if err := c.PersistentFlags().parse(args[:idx]); err != nil {
							if errors.Is(err, flag.ErrHelp) {
								fmt.Println(c.help())
								return nil
							}
							return err
						}
					}
					parsedPersistent := []string{}
					if c.persistentFlags != nil {
						parsedPersistent = c.PersistentFlags().args()
					}
					if c.PersistentPreRunE != nil {
						if err := c.PersistentPreRunE(c, parsedPersistent); err != nil {
							return err
						}
					}
					sub.ctx = c.ctx
					if err := sub.execute(args[idx+1:]); err != nil {
						return err
					}
					return nil
				}
			}
		}
	}
	if c.persistentFlags != nil && c.persistentFlags.hasAny() {
		if err := c.PersistentFlags().parse(args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				fmt.Println(c.help())
				return nil
			}
			return err
		}
	}
	parsedArgs := args
	if c.persistentFlags != nil {
		parsedArgs = c.PersistentFlags().args()
	}
	if c.PersistentPreRunE != nil {
		if err := c.PersistentPreRunE(c, parsedArgs); err != nil {
			return err
		}
	}

	if err := c.Flags().parse(parsedArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Println(c.help())
			return nil
		}
		return err
	}
	parsedArgs = c.Flags().args()

	if err := c.Flags().validateRequired(); err != nil {
		return err
	}

	if c.RunE != nil {
		return c.RunE(c, parsedArgs)
	}
	if c.Run != nil {
		c.Run(c, parsedArgs)
		return nil
	}

	if len(c.subcommands) > 0 {
		fmt.Println(c.help())
		return nil
	}
	return errors.New("no command handler")
}

func (c *Command) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *Command) Flags() *FlagSet {
	if c.flags == nil {
		c.flags = newFlagSet(c.Use)
	}
	return c.flags
}

func (c *Command) PersistentFlags() *FlagSet {
	if c.persistentFlags == nil {
		c.persistentFlags = newFlagSet(c.Use + "-persistent")
	}
	return c.persistentFlags
}

func (c *Command) MarkFlagRequired(name string) error {
	c.Flags().required[name] = struct{}{}
	return nil
}

func (c *Command) help() string {
	if len(c.subcommands) == 0 {
		return c.Use
	}
	lines := []string{"Usage:", "  " + c.Use, "", "Available Commands:"}
	for _, sub := range c.subcommands {
		lines = append(lines, fmt.Sprintf("  %-12s %s", strings.Fields(sub.Use)[0], sub.Short))
	}
	return strings.Join(lines, "\n")
}

type FlagSet struct {
	fs       *flag.FlagSet
	required map[string]struct{}
}

func newFlagSet(name string) *FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	return &FlagSet{fs: fs, required: map[string]struct{}{}}
}

func (f *FlagSet) StringVar(p *string, name, value, usage string) {
	f.fs.StringVar(p, name, value, usage)
}

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.fs.BoolVar(p, name, value, usage)
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.fs.IntVar(p, name, value, usage)
}

func (f *FlagSet) StringSliceVar(p *[]string, name string, value []string, usage string) {
	joined := ""
	if len(value) > 0 {
		joined = strings.Join(value, ",")
	}
	f.fs.Func(name, usage, func(v string) error {
		if v == "" {
			*p = nil
			return nil
		}
		*p = strings.Split(v, ",")
		return nil
	})
	if joined != "" {
		_ = f.fs.Set(name, joined)
	}
}

func (f *FlagSet) parse(args []string) error {
	return f.fs.Parse(args)
}

func (f *FlagSet) args() []string {
	return f.fs.Args()
}

func (f *FlagSet) validateRequired() error {
	missing := make([]string, 0)
	f.fs.VisitAll(func(fl *flag.Flag) {
		if _, ok := f.required[fl.Name]; !ok {
			return
		}
		if fl.Value.String() == "" {
			missing = append(missing, fl.Name)
		}
	})
	if len(missing) > 0 {
		return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
	}
	return nil
}

func (f *FlagSet) hasAny() bool {
	has := false
	f.fs.VisitAll(func(_ *flag.Flag) {
		has = true
	})
	return has
}

func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" || a == "-help" || a == "help" {
			return true
		}
	}
	return false
}
