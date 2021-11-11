package main

import (
	"flag"
	"fmt"
	interact "github.com/hyperjumptech/hyper-interactive"
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

	flag.Parse()

	configFile := *configFilePtr
	startMinion := *minionPtr
	startCentral := *centralPtr
	runOnce := *runOncePtr
	setup := *setupPtr

	if setup {
		Setup(configFile)
	} else if startMinion {
		StartMinion(configFile)
	} else if startCentral {
		StartServer(configFile)
	} else if len(runOnce) > 0 {
		ProbeOnce(runOnce, configFile)
	} else {
		flag.Usage = func() {
			fmt.Fprintf(flag.CommandLine.Output(), "Usage : %s (-central|-minion|-config|-once <probe>|-setup) -config <config-file>\n  Arguments:\n", os.Args[0])
			flag.PrintDefaults()
		}
		flag.Usage()
	}
}

func Setup(config string) {
	if len(config) == 0 {
		if fInfo, err := os.Stat("./mihp.yaml"); err == nil && !fInfo.IsDir() {
			SetupConfig("./mihp.yaml")
		} else if fInfo, err := os.Stat("/etc/mihp/mihp.yaml"); err == nil && !fInfo.IsDir() {
			SetupConfig("/etc/mihp/mihp.yaml")
		} else {
			if interact.Confirm("You do not specify configuration file to configure, you want to create one in current folder \"./mihp.yaml\" ? ", true) {
				SetupConfig("./mihp.yaml")
			}
		}
	} else {
		if fInfo, err := os.Stat(config); err != nil {
			fmt.Printf("Problem open file %s, got %s\n", config, err.Error())
			if interact.Confirm("Do you want to create one in this directory \"./mihp.yaml\" ? ", true) {
				SetupConfig("./mihp.yaml")
			}
		} else if fInfo.IsDir() {
			fmt.Printf("Problem open file %s, its a directory", config)
			if interact.Confirm(fmt.Sprintf("Do you want to create one in ? \"%s/mihp.yaml\" ? ", config), true) {
				SetupConfig(fmt.Sprintf("%s/mihp.yaml", config))
			}
		} else {
			SetupConfig(config)
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

}
