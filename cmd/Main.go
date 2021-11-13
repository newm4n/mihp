package main

import (
	"context"
	"flag"
	"fmt"
	interact "github.com/hyperjumptech/hyper-interactive"
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/minion/probing"
	"io/ioutil"
	"os"
)

var (
	splash = `
   _____  .___  ___ _____________ 
  /     \ |   |/   |   \______   \
 /  \ /  \|   /    ~    \     ___/
/    Y    \   \    Y    /    |    
\____|__  /___|\___|_  /|____|    
        \/           \/           
       MIHP Is HTTP Probe
`
)

func main() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s\n", splash)

	minionPtr := flag.Bool("minion", false, "Start probe as minion / probe node")
	centralPtr := flag.Bool("central", false, "Start central probe management server")
	runOncePtr := flag.String("once", "", "Probe name to run once when minion is started. Use in conjunction with -minion. Probe result will displayed directly in the console")
	setupPtr := flag.Bool("setup", false, "Create/Modify a configuration file interactively")
	configFilePtr := flag.String("config", "", "Configuration file to use.")
	helpPtr := flag.Bool("help", false, "Show this help.")

	flag.Parse()

	configFile := *configFilePtr
	startMinion := *minionPtr
	startCentral := *centralPtr
	runOnce := *runOncePtr
	setup := *setupPtr
	help := *helpPtr

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage : %s (-central|-minion|-config|-once <probe>|-setup) -config <config-file>\n  Arguments:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if help {
		flag.Usage()
	} else if setup {
		Setup(configFile)
	} else if len(runOnce) > 0 {
		ProbeOnce(runOnce, configFile)
	} else if startMinion {
		StartMinion(configFile)
	} else if startCentral {
		StartServer(configFile)
	} else {
		flag.Usage()
	}
}

func Setup(config string) {
	if len(config) == 0 {
		if fInfo, err := os.Stat("./mihp.yaml"); err == nil && !fInfo.IsDir() {
			err := SetupConfig("./mihp.yaml")
			if err != nil {
				fmt.Printf("got error %s\n", err.Error())
			}
		} else if fInfo, err := os.Stat("/etc/mihp/mihp.yaml"); err == nil && !fInfo.IsDir() {
			err := SetupConfig("/etc/mihp/mihp.yaml")
			if err != nil {
				fmt.Printf("got error %s\n", err.Error())
			}
		} else {
			if interact.Confirm("You do not specify configuration file to configure, you want to create one in current folder \"./mihp.yaml\" ? ", true) {
				err := SetupConfig("./mihp.yaml")
				if err != nil {
					fmt.Printf("got error %s\n", err.Error())
				}
			}
		}
	} else {
		if fInfo, err := os.Stat(config); err != nil {
			fmt.Printf("Problem open file %s, got %s\n", config, err.Error())
			if interact.Confirm("Do you want to create one in this directory \"./mihp.yaml\" ? ", true) {
				err := SetupConfig("./mihp.yaml")
				if err != nil {
					fmt.Printf("got error %s\n", err.Error())
				}
			}
		} else if fInfo.IsDir() {
			fmt.Printf("Problem open file %s, its a directory", config)
			if interact.Confirm(fmt.Sprintf("Do you want to create one in ? \"%s/mihp.yaml\" ? ", config), true) {
				err := SetupConfig(fmt.Sprintf("%s/mihp.yaml", config))
				if err != nil {
					fmt.Printf("got error %s\n", err.Error())
				}
			}
		} else {
			err := SetupConfig(config)
			if err != nil {
				fmt.Printf("got error %s\n", err.Error())
			}
		}
	}
	fmt.Println("Bye.")
}

func StartServer(config string) {
	fmt.Println("Bye.")
}

func StartMinion(config string) {
	fmt.Println("Bye.")
}

func ProbeOnce(probeName, config string) {
	fmt.Println("Bye.")
	file, err := os.Open(config)
	if err != nil {
		fmt.Printf("Got error while opening %s got %s", config, err.Error())
		return
	}
	yamlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Got error while reading %s got %s", config, err.Error())
		return
	}
	cfg, err := internal.YAMLToMIHPConfig(yamlBytes)
	if err != nil {
		fmt.Printf("Got error while parsing yaml content in %s got %s", config, err.Error())
		return
	}
	if cfg.Version != internal.Version {
		fmt.Printf("Wrong YAML version in file %s got %s instead of %s", config, cfg.Version, internal.Version)
		return
	}
	if cfg.ProbePool == nil {
		fmt.Printf("Configuration file %s contains no probe", config)
		return
	}
	for _, probe := range cfg.ProbePool {
		if probe.Name == probeName {
			timeout := 10
			pCtx := internal.NewProbeContext()
			fmt.Printf("Probing once. time-out %d seconds", timeout)
			err := probing.ExecuteProbe(context.Background(), probe, pCtx, timeout, true, true)
			if err != nil {
				fmt.Printf("Error while runing probe %s. Got %s", probeName, err.Error())
				return
			}
			fmt.Printf("Successfuly execute probe %s. Context follows.\n%s\n", probeName, pCtx.ToString(false))
			return
		}
	}

	fmt.Printf("Configuration file %s do not contain probe named %s.", config, probeName)
}
