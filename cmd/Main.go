package main

import (
	"context"
	"flag"
	"fmt"
	hyper_interactive "github.com/hyperjumptech/hyper-interactive"
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/internal/probing"
	"github.com/newm4n/mihp/minion"
	"github.com/newm4n/mihp/pkg/errors"
	"github.com/newm4n/mihp/pkg/helper/cron"
	"io/ioutil"
	"net/url"
	"os"
	"time"
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
		fmt.Fprintf(flag.CommandLine.Output(), "Visit https://github.com/newm4n/mihp/documentation.md to know how to use MIHP\n")
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

func LoadConfigFile(configPath string) (config *internal.MIHPConfig, err error) {
	if len(configPath) == 0 {
		if fInfo, err := os.Stat("./mihp.yaml"); err == nil && !fInfo.IsDir() {
			return loadConfig("./mihp.yaml")
		}
		if fInfo, err := os.Stat("/etc/mihp/mihp.yaml"); err == nil && !fInfo.IsDir() {
			return loadConfig("/etc/mihp/mihp.yaml")
		}
		return nil, errors.ErrConfigFileNotFound
	}
	if fInfo, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("can not load config file %s. got %w", configPath, err)
	} else if fInfo.IsDir() {
		return nil, fmt.Errorf("can not load config file %s, its a directory", configPath)
	}
	return loadConfig(configPath)
}

func loadConfig(path string) (config *internal.MIHPConfig, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	yamlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	cfg, err := internal.YAMLToMIHPConfig(yamlBytes)
	if err != nil {
		return nil, err
	}
	if cfg.Version != internal.Version {
		return nil, fmt.Errorf("invalid YAML Version %s", cfg.Version)
	}
	if cfg.ProbePool == nil {
		cfg.ProbePool = make(internal.ProbePool, 0)
	}
	return cfg, nil
}

func Setup(config string) {
	cfg, err := LoadConfigFile(config)
	if err != nil {
		if err == errors.ErrConfigFileNotFound {
			if hyper_interactive.Confirm("MIHP.yaml not found, you wish to create one in the current folder ?", true) {
				err := SetupConfig(&internal.MIHPConfig{
					Version:   internal.Version,
					ProbePool: make(internal.ProbePool, 0),
				}, "./MIHP.yaml")
				if err != nil {
					fmt.Printf("got error %s\n", err.Error())
				}
			}
		} else {
			fmt.Printf("got error %s\n", err.Error())
		}
	} else {
		err := SetupConfig(cfg, config)
		if err != nil {
			fmt.Printf("got error %s\n", err.Error())
		}
	}
	fmt.Println("Bye.")
}

func StartServer(config string) {
	fmt.Println("Bye.")
}

func StartMinion(config string) {
	cfg, err := LoadConfigFile(config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "got error %s\n", err.Error())
		return
	} else {
		if cfg.Minion == nil {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion section\n")
			return
		}
		if len(cfg.Minion.MinionUID) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion UID\n")
			return
		}
		if len(cfg.Minion.Name) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion Name\n")
			return
		}
		if len(cfg.Minion.CentralBaseURL) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion Central URL\n")
			return
		} else {
			_, err := url.ParseRequestURI(cfg.Minion.CentralBaseURL)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "invalid minions central URL %s\n", cfg.Minion.CentralBaseURL)
				return
			}
		}
		if len(cfg.Minion.Datacenter) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion Data Center name\n")
			return
		}
		if len(cfg.Minion.CountryISO) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion Country ISO code\n")
			return
		} else {
			found := false
			for _, cc := range CountryCodes {
				if cc.Code == cfg.Minion.CountryISO {
					found = true
					break
				}
			}
			if !found {
				_, _ = fmt.Fprintf(os.Stderr, "invalid minion country code %s\n", cfg.Minion.CountryISO)
				return
			}
		}
		if len(cfg.Minion.ReportCron) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "configuration missing minion Reporting CRON schedule\n")
			return
		} else {
			_, err = cron.NewSchedule(cfg.Minion.ReportCron)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "invalid minion cron syntax %s\n", cfg.Minion.ReportCron)
				return
			}
		}

		fmt.Println("Starting MIHP MINION")
		minionContext := context.Background()
		minion.Start(minionContext, cfg)
	}
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
			fmt.Printf("Probing once. time-out %d seconds\n", timeout)
			err := probing.ExecuteProbe(context.Background(), probe, pCtx, timeout, true, true)
			if err != nil {
				fmt.Printf("Error while runing probe %s. Got %s.\n", probe.Name, err.Error())
				fmt.Printf("If the error is about I/O, check for some firewall or VPN.\n")
				path := fmt.Sprintf("./probe-%s-fail-%s.txt", probe.Name, time.Now().Format(time.RFC3339))
				file, err := os.Create(path)
				if err != nil {
					fmt.Printf("can not write context dump to %s\n", path)
					return
				}
				defer file.Close()
				file.WriteString(pCtx.ToString(false))
				fmt.Printf("Context written to %s\n", path)
				os.Exit(1)
			}
			fmt.Printf("Successfuly execute probe %s.\n", probe.Name)
			path := fmt.Sprintf("./probe-%s-success-%s.txt", probe.Name, time.Now().Format(time.RFC3339))
			file, err := os.Create(path)
			if err != nil {
				fmt.Printf("can not write context dump to %s\n", path)
				return
			}
			defer file.Close()
			file.WriteString(pCtx.ToString(false))
			fmt.Printf("Context written to %s\n", path)
			os.Exit(0)
		}
	}

	fmt.Printf("Configuration file %s do not contain probe named %s.", config, probeName)
}
