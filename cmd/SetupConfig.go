package main

import (
	"fmt"
	"github.com/google/uuid"
	interact "github.com/hyperjumptech/hyper-interactive"
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/pkg/helper/cron"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"strings"
)

func SetupConfig(configFile string) (err error) {
	//Mode := ""
	var Config *internal.MIHPConfig
	fInfo, err := os.Stat(configFile)
	if err != nil || fInfo.IsDir() {
		//Mode = "NEW"
		Config = &internal.MIHPConfig{}
	} else {
		file, err := os.Open(configFile)
		if err != nil {
			return err
		}
		yamlBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		cfg, err := internal.YAMLToMIHPConfig(yamlBytes)
		if err != nil {
			return err
		}
		//Mode = "EDIT"
		Config = cfg
	}
	return showMainMenu(Config, configFile)
}

func showMainMenu(config *internal.MIHPConfig, configFile string) (err error) {
	for {
		fmt.Println("\n---[ MAIN CONFIGURATION ]-------------------------------------")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "STATUS"})

		if config.ProbePool == nil {
			table.Append([]string{"Probes", "Not Configured"})
		} else {
			table.Append([]string{"Probes", fmt.Sprintf("%d configured", len(config.ProbePool))})
		}
		if config.Central == nil {
			table.Append([]string{"Central", "Not Configured"})
		} else {
			table.Append([]string{"Central", "Configured"})
		}
		if config.Minion == nil {
			table.Append([]string{"Minion", "Not Configured"})
		} else {
			table.Append([]string{"Minion", "Configured"})
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{"Manage probes", "Configure Central", "Configure Minion", "Save and Finish"}, 1, 3, false)

		switch selected {
		case 1:
			err = manageProbe(config)
			if err != nil {
				return err
			}
		case 2:
			err = configureCentral(config)
			if err != nil {
				return err
			}
		case 3:
			err = configureMinion(config)
			if err != nil {
				return err
			}
		case 4:
			err = saveAndExist(config, configFile)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

func manageProbe(config *internal.MIHPConfig) (err error) {
	return nil
}

func configureCentral(config *internal.MIHPConfig) (err error) {
	central := config.Central
	if central == nil {
		central = &internal.CentralConfig{}
	}
	for {
		fmt.Println("\n---[ CENTRAL CONFIGURATION ]-------------------------------------")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "STATUS"})

		if len(central.ListenHost) == 0 {
			table.Append([]string{"Server Host", "Not Configured"})
		} else {
			table.Append([]string{"Server Host", central.ListenHost})
		}
		if central.ListenPort == 0 {
			table.Append([]string{"Server Port", "Not Configured"})
		} else {
			table.Append([]string{"Server Port", fmt.Sprintf("%d", central.ListenPort)})
		}
		if central.ReadHeaderTimeoutSecond == 0 {
			table.Append([]string{"Header Timeout", "Not Configured"})
		} else {
			table.Append([]string{"Header Timeout", fmt.Sprintf("%d second", central.ReadHeaderTimeoutSecond)})
		}
		if central.ReadTimeoutSecond == 0 {
			table.Append([]string{"Read Timeout", "Not Configured"})
		} else {
			table.Append([]string{"Read Timeout", fmt.Sprintf("%d second", central.ReadTimeoutSecond)})
		}
		if central.WriteTimeoutSecond == 0 {
			table.Append([]string{"Write Timeout", "Not Configured"})
		} else {
			table.Append([]string{"Write Timeout", fmt.Sprintf("%d second", central.WriteTimeoutSecond)})
		}
		if central.IdleTimeoutSecond == 0 {
			table.Append([]string{"Idle Timeout", "Not Configured"})
		} else {
			table.Append([]string{"Idle Timeout", fmt.Sprintf("%d second", central.IdleTimeoutSecond)})
		}
		if len(central.AdminUser) == 0 {
			table.Append([]string{"Admin User", "Not Configured"})
		} else {
			table.Append([]string{"Admin User", central.AdminUser})
		}
		if len(central.AdminPassword) == 0 {
			table.Append([]string{"Admin Password", "Not Configured"})
		} else {
			table.Append([]string{"Admin Password", central.AdminPassword})
		}
		if len(central.JWTIssuer) == 0 {
			table.Append([]string{"JWT Issuer", "Not Configured"})
		} else {
			table.Append([]string{"JWT Issuer", central.JWTIssuer})
		}
		if len(central.JWTKey) == 0 {
			table.Append([]string{"JWT Key", "Not Configured"})
		} else {
			table.Append([]string{"JWT Key", central.JWTKey})
		}
		if central.JWTAccessKeyAgeMinute == 0 {
			table.Append([]string{"JWT Access Key Age", "Not Configured"})
		} else {
			table.Append([]string{"JWT Access Key Age", fmt.Sprintf("%d minutes", central.JWTAccessKeyAgeMinute)})
		}
		if central.JWTRefreshKeyAgeMinute == 0 {
			table.Append([]string{"JWT Refresh Key Age", "Not Configured"})
		} else {
			table.Append([]string{"JWT Refresh Key Age", fmt.Sprintf("%d minutes", central.JWTRefreshKeyAgeMinute)})
		}

		if central.MySQLConfig == nil {
			table.Append([]string{"MySQL", "Not Configured"})
		} else {
			table.Append([]string{"MySQL", "Configured"})
		}
		if central.PostgreSQLConfig == nil {
			table.Append([]string{"PostgreSQL", "Not Configured"})
		} else {
			table.Append([]string{"PostgreSQL", "Configured"})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Set Server Host", "Set Server Port",
			"Set Header Read Timeout", "Set Read Timout",
			"Set Write Timeout", "Set Idle Timout",
			"Set Admin User", "Set Admin Password",
			"Set JWT Issuer", "Set JWT Key",
			"Set JWT Access Token Age", "Set JWT Refresh Token Age",
			"Configure MySQL", "Configure PostgreSQL",
			"Load Defaults", "Finish"}, 1, 3, false)
		switch selected {
		case 1:
			var defa string
			if len(central.ListenHost) > 0 {
				defa = central.ListenHost
			} else {
				defa = "0.0.0.0"
			}
			central.ListenHost = interact.Ask("Specify New Server Host IP", defa, true)
		case 2:
			var defa int
			if central.ListenPort == 0 {
				defa = 53421
			} else {
				defa = central.ListenPort
			}
			central.ListenPort = interact.AskNumber("Pick New Port Number", 1025, 65000, defa, true)
		case 3:
			var defa int
			if central.ReadHeaderTimeoutSecond == 0 {
				defa = 10
			} else {
				defa = central.ReadHeaderTimeoutSecond
			}
			central.ReadHeaderTimeoutSecond = interact.AskNumber("Specify New Read-Header-Timeout in Second", 3, 120, defa, true)
		case 4:
			var defa int
			if central.ReadTimeoutSecond == 0 {
				defa = 10
			} else {
				defa = central.ReadTimeoutSecond
			}
			central.ReadTimeoutSecond = interact.AskNumber("Specify New Read-Timeout in Second", 3, 120, defa, true)
		case 5:
			var defa int
			if central.WriteTimeoutSecond == 0 {
				defa = 10
			} else {
				defa = central.WriteTimeoutSecond
			}
			central.WriteTimeoutSecond = interact.AskNumber("Specify New Write-Timeout in Second", 3, 120, defa, true)
		case 6:
			var defa int
			if central.IdleTimeoutSecond == 0 {
				defa = 10
			} else {
				defa = central.IdleTimeoutSecond
			}
			central.IdleTimeoutSecond = interact.AskNumber("Specify New Idle-Timeout in Second", 3, 120, defa, true)
		case 7:
			var defa string
			if len(central.AdminUser) > 0 {
				defa = central.AdminUser
			} else {
				defa = "Admin"
			}
			central.AdminUser = interact.Ask("Specify New Admin account name", defa, true)
		case 8:
			var defa string
			if len(central.AdminPassword) > 0 {
				defa = central.AdminPassword
			} else {
				defa = "This is a very good pass Phrase"
			}
			central.AdminPassword = interact.Ask("Specify New Admin Password", defa, true)
		case 9:
			var defa string
			if len(central.JWTIssuer) > 0 {
				defa = central.JWTIssuer
			} else {
				defa = "mihp.io"
			}
			central.JWTIssuer = interact.Ask("Specify New JWT Issuer", defa, true)
		case 10:
			var defa string
			if len(central.JWTKey) > 0 {
				defa = central.JWTKey
			} else {
				defa = "Th1sk3ymu$tb3ch4ngebef0r3use@ge0npr0duction"
			}
			central.JWTKey = interact.Ask("Specify New JWT Secret Key", defa, true)
		case 11:
			var defa int
			if central.JWTAccessKeyAgeMinute == 0 {
				defa = 10
			} else {
				defa = central.JWTAccessKeyAgeMinute
			}
			central.JWTAccessKeyAgeMinute = interact.AskNumber("Specify New Access Token Age in Minute", 3, 60*24*365*10, defa, true)
		case 12:
			var defa int
			if central.JWTRefreshKeyAgeMinute == 0 {
				defa = 60 * 24 * 365 * 2
			} else {
				defa = central.JWTRefreshKeyAgeMinute
			}
			central.JWTRefreshKeyAgeMinute = interact.AskNumber("Specify New Refresh Token Age in Minute", 3, 60*24*365*10, defa, true)
		case 13:
			newMySQLDBConfig := ConfigureDatabase("MySQL", central.MySQLConfig)
			central.MySQLConfig = newMySQLDBConfig
		case 14:
			newPostgreSQLDBConfig := ConfigureDatabase("PostgreSQL", central.PostgreSQLConfig)
			central.PostgreSQLConfig = newPostgreSQLDBConfig
		case 15:
			central.AdminUser = "Admin"
			central.AdminPassword = "super secret pass phrase"
			central.JWTIssuer = "mihp.io"
			central.JWTKey = "45jhew3?54ej3u%(3puo2^$3f433$<k4jh34908rkjccmoie4n43mmnv$k"
			central.JWTRefreshKeyAgeMinute = 365 * 2
			central.JWTAccessKeyAgeMinute = 10
			central.ListenHost = "0.0.0.0"
			central.ListenPort = 53829
			central.IdleTimeoutSecond = 30
			central.ReadHeaderTimeoutSecond = 5
			central.ReadTimeoutSecond = 10
			central.WriteTimeoutSecond = 50
			central.MySQLConfig = &internal.DBConfig{
				Host:     "0.0.0.0",
				Port:     3306,
				User:     "sa",
				Password: "sa",
				Database: "mihp",
			}
			central.PostgreSQLConfig = &internal.DBConfig{
				Host:     "0.0.0.0",
				Port:     5432,
				User:     "sa",
				Password: "sa",
				Database: "mihp",
			}
		case 16:
			if config.Central == nil {
				config.Central = central
			}
			return nil
		}

	}

	return nil
}

func ConfigureDatabase(dbName string, cfg *internal.DBConfig) *internal.DBConfig {
	if cfg == nil {
		cfg = &internal.DBConfig{}
	}
	for {
		fmt.Printf("\n---[ %s DATABASE CONFIGURATION ]-------------------------------------\n", dbName)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "STATUS"})

		if len(cfg.Host) == 0 {
			table.Append([]string{"DB Host", "Not Configured"})
		} else {
			table.Append([]string{"DB Host", cfg.Host})
		}
		if cfg.Port == 0 {
			table.Append([]string{"DB Port", "Not Configured"})
		} else {
			table.Append([]string{"DB Port", fmt.Sprintf("%d", cfg.Port)})
		}

		if len(cfg.User) == 0 {
			table.Append([]string{"DB User", "Not Configured"})
		} else {
			table.Append([]string{"DB User", cfg.User})
		}
		if len(cfg.Password) == 0 {
			table.Append([]string{"DB Password", "Not Configured"})
		} else {
			table.Append([]string{"DB Password", cfg.Password})
		}
		if len(cfg.Database) == 0 {
			table.Append([]string{"DB schema", "Not Configured"})
		} else {
			table.Append([]string{"DB schema", cfg.Database})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		switch interact.Select("What to do ?", []string{
			"Set DB Host IP", "Set DB Port", "Set DB User", "Set DB Password",
			"Set DB Schema", "Finish",
		}, 1, 6, false) {
		case 1:
			var defa string
			if len(cfg.Host) > 0 {
				defa = cfg.Host
			} else {
				defa = "0.0.0.0"
			}
			cfg.Host = interact.Ask("Specify New DB Host IP", defa, true)
		case 2:
			var defa int
			if cfg.Port != 0 {
				defa = cfg.Port
			} else {
				if strings.ToUpper(dbName) == "MYSQL" {
					defa = 3306
				} else if strings.ToUpper(dbName) == "POSTGRESQL" {
					defa = 5432
				} else {
					defa = 4321
				}
			}
			cfg.Port = interact.AskNumber("Specify new port number", 1024, 65000, defa, true)
		case 3:
			var defa string
			if len(cfg.User) > 0 {
				defa = cfg.User
			} else {
				defa = "root"
			}
			cfg.User = interact.Ask("Specify New DB User", defa, true)
		case 4:
			var defa string
			if len(cfg.Password) > 0 {
				defa = cfg.Password
			} else {
				defa = "this is db user password"
			}
			cfg.Password = interact.Ask("Specify New DB Password", defa, true)
		case 5:
			var defa string
			if len(cfg.Database) > 0 {
				defa = cfg.Database
			} else {
				defa = "this is db user password"
			}
			cfg.Database = interact.Ask("Specify New DB Schema", defa, true)
		case 6:
			return cfg
		}
	}
}

func configureMinion(config *internal.MIHPConfig) (err error) {
	minion := config.Minion
	if minion == nil {
		minion = &internal.MinionConfig{}
	}

	for {
		fmt.Println("\n---[ MINION CONFIGURATION ]-------------------------------------")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "STATUS"})

		if len(minion.Name) == 0 {
			table.Append([]string{"Name", "Not Configured"})
		} else {
			table.Append([]string{"Name", minion.Name})
		}
		if len(minion.MinionUID) == 0 {
			table.Append([]string{"UID", "Not Configured"})
		} else {
			table.Append([]string{"UID", minion.MinionUID})
		}
		if len(minion.CountryISO) == 0 {
			table.Append([]string{"Country Code", "Not Configured"})
		} else {
			table.Append([]string{"Country Code", minion.CountryISO})
		}
		if len(minion.Datacenter) == 0 {
			table.Append([]string{"Data Center", "Not Configured"})
		} else {
			table.Append([]string{"Data Center", minion.Datacenter})
		}
		if len(minion.CentralBaseURL) == 0 {
			table.Append([]string{"Central Base URL", "Not Configured"})
		} else {
			table.Append([]string{"Central Base URL", minion.CentralBaseURL})
		}
		if len(minion.ReportCron) == 0 {
			table.Append([]string{"Reporting CRON", "Not Configured"})
		} else {
			table.Append([]string{"Reporting CRON", minion.ReportCron})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		switch interact.Select("What to do ?", []string{
			"Set Name", "Set UID", "Set Country", "Set Datacenter",
			"Set Central Base URL", "Set Reporting Cron",
			"Finish",
		}, 1, 7, false) {
		case 1:
			minion.CentralBaseURL = interact.Ask("New Name ?", "MyMinion", true)
		case 2:
			if interact.Confirm("You want us to autogenerate UID for you ?", true) {
				minion.MinionUID = uuid.New().String()
			} else {
				minion.MinionUID = interact.Ask("New UID for this minion ?", uuid.New().String(), false)
			}
		case 3:
			for {
				if interact.Confirm("You know your country code ?", true) {
					cc := interact.Ask("What is the country code ?", "US", false)
					if len(CountryNameForCode(cc)) > 0 {
						minion.CountryISO = strings.ToUpper(cc)
						break
					} else {
						continue
					}
				} else {
					ccidx := interact.Select("Choose your country code", CountryCodeOptions(), 0, 10, true)
					minion.CountryISO = CountryCodes[ccidx].Code
					break
				}
			}
		case 5:
			minion.CentralBaseURL = interact.Ask("New Central Base URL ?", "https://hyperjump.tech", true)
		case 6:
			for {
				minion.ReportCron = interact.Ask("New Reporting CRON ?", "0 */5 * * * * *", true)
				if _, err := cron.NewSchedule(minion.ReportCron); err != nil {
					fmt.Println("Invalid CRON syntax")
					continue
				} else if strings.HasPrefix(minion.ReportCron, "* ") {
					fmt.Println("Reporting interval to short")
					continue
				} else {
					break
				}
			}
		case 7:
			if config.Minion == nil {
				config.Minion = minion
			}
			return err
		}
	}
}

func saveAndExist(config *internal.MIHPConfig, configFile string) (err error) {
	fmt.Printf("Saving to %s\n", configFile)
	// todo finish this saving
	return nil
}
