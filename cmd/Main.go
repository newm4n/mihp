package main

import (
	"flag"
	"fmt"
)

var (
	splash = `   _____  .___  ___ _____________ 
  /     \ |   |/   |   \______   \
 /  \ /  \|   /    ~    \     ___/
/    Y    \   \    Y    /    |    
\____|__  /___|\___|_  /|____|    
        \/           \/           
       MIHP Is HTTP Probe`
)

func main() {
	fmt.Println(splash)

	minionPtr := flag.Bool("minion", false, "Start probe as minion / probe node")
	centralPtr := flag.Bool("central", false, "Start central probe management server")
	configFilePtr := flag.String("config", "", "Configuration file to use.")
	runOncePtr := flag.String("once", "", "Probe name to run once when minion is started. Use in conjunction with -minion. Probe result will displayed directly in the console")

	flag.Parse()

	configFile := *configFilePtr
	startMinion := *minionPtr
	startCentral := *centralPtr
	runOnce := *runOncePtr

}
