package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/marinator86/portier-cli/internal/portier/application"
	"github.com/marinator86/portier-cli/internal/utils"
	"github.com/spf13/cobra"
)

type runOptions struct {
	ConfigFile   string
	ApiTokenFile string
	Output       string
}

func defaultRunOptions() (*runOptions, error) {
	home, err := utils.Home()
	if err != nil {
		log.Printf("could not get home directory: %v", err)
		return nil, err
	}

	configFile := home + "/config.yaml"
	apiTokenFile := home + "/credentials_device.yaml"

	return &runOptions{
		ConfigFile:   configFile,
		ApiTokenFile: apiTokenFile,
		Output:       "json",
	}, nil
}

func newRunCmd() (*cobra.Command, error) {
	o, err := defaultRunOptions()
	if err != nil {
		log.Printf("could not get default options: %v", err)
		return nil, err
	}

	cmd := &cobra.Command{
		Use:          "run",
		Short:        "run the relay, i.e. all local services defined for this device (requires registration)",
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
		RunE:         o.run,
	}

	cmd.Flags().StringVarP(&o.ConfigFile, "config file", "c", o.ConfigFile, "config file path")
	cmd.Flags().StringVarP(&o.ApiTokenFile, "apiToken file", "t", o.ApiTokenFile, "apiToken file path")
	cmd.Flags().StringVarP(&o.Output, "output", "o", o.Output, "output format (yaml | json)")

	return cmd, nil
}

func (o *runOptions) run(cmd *cobra.Command, args []string) error {
	err := o.parseArgs(cmd, args)
	if err != nil {
		log.Println("could not parse args")
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "starting device, services %s, apiToken %s, out %s\n", o.ConfigFile, o.ApiTokenFile, o.Output)

	application := application.NewPortierApplication()

	application.LoadConfig(o.ConfigFile)

	application.LoadApiToken(o.ApiTokenFile)

	application.StartServices()

	// wait until process is killed
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	application.StopServices()

	return nil
}

func (o *runOptions) parseArgs(cmd *cobra.Command, _ []string) error {
	configFile, err := cmd.Flags().GetString("config")
	if err == nil {
		o.ConfigFile = configFile
	}

	apiTokenFile, err := cmd.Flags().GetString("apiToken")
	if err == nil {
		o.ApiTokenFile = apiTokenFile
	}

	output, err := cmd.Flags().GetString("output")
	if err == nil {
		o.Output = output
	}

	return nil
}
