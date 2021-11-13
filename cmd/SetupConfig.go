package main

import (
	"fmt"
	"github.com/google/uuid"
	interact "github.com/hyperjumptech/hyper-interactive"
	"github.com/newm4n/mihp/internal"
	"github.com/newm4n/mihp/pkg/helper"
	"github.com/newm4n/mihp/pkg/helper/cron"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
	"regexp"
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
	if config.ProbePool == nil {
		config.ProbePool = make(internal.ProbePool, 0)
	}
	for {
		fmt.Println("\n---[ PROBE LIST ]-------------------------------------")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NO", "PROBE NAME", "BASE URL", "CRON", "REQUESTS"})
		for idx, prob := range config.ProbePool {
			table.Append([]string{fmt.Sprintf("%d", idx), prob.Name, prob.BaseURL, prob.Cron, fmt.Sprintf("%d", len(prob.Requests))})
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Add Probe", "Remove Probe", "View/Update Probe", "Finish"}, 1, 4, false)

		probeList := make([]string, len(config.ProbePool))
		for idx, p := range config.ProbePool {
			probeList[idx] = p.Name
		}

		switch selected {
		case 1:
			addProbe(config.ProbePool)
		case 2:
			if len(probeList) == 0 {
				fmt.Println("There are no probe to remove.")
			} else {
				probeList := append(probeList, "Cancel")
				sel := interact.Select("Select probe to remove", probeList, 0, len(probeList)-1, false)
				if sel == len(probeList)-1 {
					continue
				}
				if interact.Confirm(fmt.Sprintf("Are you sure to remove probe %s", config.ProbePool[sel].Name), false) {
					config.ProbePool = append(config.ProbePool[:sel], config.ProbePool[sel+1:]...)
					probeList = append(probeList[:sel], probeList[sel+1:]...)
					return
				} else {
					continue
				}
			}
		case 3:
			if len(probeList) == 0 {
				fmt.Println("There are no probe to view/edit.")
			} else {
				probeList := append(probeList, "Cancel")
				sel := interact.Select("Select probe to view/edit", probeList, 0, len(probeList)-1, false)
				if sel == len(probeList)-1 {
					continue
				}
				editProbe(config.ProbePool, sel)
			}
		case 4:
			return nil
		}
	}
}

var (
	spaceChecker = regexp.MustCompile(`\s`)
)

func askNoSpaceNotEmpty(question, fieldName, defa string, confirm bool) string {
	for {
		answer := interact.Ask(question, defa, confirm)
		if len(answer) == 0 {
			fmt.Printf("%s must not empty", fieldName)
			continue
		}
		if spaceChecker.MatchString(answer) {
			fmt.Printf("%s must not contains space", fieldName)
			continue
		}
		return answer
	}
}

func addProbe(pool internal.ProbePool) {
	var newName, newID string
	for {
		newName = askNoSpaceNotEmpty("New probe name?", "Name", helper.RandomName(), false)
		for _, p := range pool {
			if p.Name == newName {
				fmt.Println("That name already exist in list of probes")
				continue
			}
		}
		break
	}
	for {
		newID = askNoSpaceNotEmpty("New probe ID?", "ID", uuid.New().String(), false)
		for _, p := range pool {
			if p.ID == newID {
				fmt.Println("That ID already exist in list of probes")
				continue
			}
		}
		break
	}
	p := &internal.Probe{
		Name:     newName,
		ID:       newID,
		Requests: make([]*internal.ProbeRequest, 0),
	}
	pool = append(pool, p)
	editProbe(pool, len(pool)-1)
}

func editProbe(pool internal.ProbePool, idx int) {
	probe := pool[idx]

	for {
		fmt.Printf("\n---[ PROBE %s CONFIGURATION ]-------------------------------------\n", probe.Name)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "VALUE"})

		if len(probe.Name) > 0 {
			table.Append([]string{"Name", probe.Name})
		} else {
			table.Append([]string{"Name", "Not Set"})
		}
		if len(probe.ID) > 0 {
			table.Append([]string{"ID", probe.ID})
		} else {
			table.Append([]string{"ID", "Not Set"})
		}
		if len(probe.BaseURL) > 0 {
			table.Append([]string{"BaseURL", probe.BaseURL})
		} else {
			table.Append([]string{"BaseURL", "Not Set"})
		}
		if len(probe.Cron) > 0 {
			table.Append([]string{"CRON", probe.Cron})
		} else {
			table.Append([]string{"CRON", "Not Set"})
		}
		table.Append([]string{"Request", fmt.Sprintf("%d configured", len(probe.Requests))})
		if probe.UpThreshold > 0 {
			table.Append([]string{"UpThreshold", fmt.Sprintf("%d", probe.UpThreshold)})
		} else {
			table.Append([]string{"UpThreshold", "Not Set"})
		}
		if probe.DownThreshold > 0 {
			table.Append([]string{"DownThreshold", fmt.Sprintf("%d", probe.DownThreshold)})
		} else {
			table.Append([]string{"DownThreshold", "Not Set"})
		}
		if probe.SMTPNotification != nil {
			table.Append([]string{"SMTP Notification", "Configured"})
		} else {
			table.Append([]string{"SMTP Notification", "Not Configured"})
		}
		if probe.CallbackNotification != nil {
			table.Append([]string{"Callback Notification", "Configured"})
		} else {
			table.Append([]string{"Callback Notification", "Not Configured"})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Set Probe Name", "Set Probe ID", "Manage Probe Requests",
			"Set Probe Base URL", "Set Probe CRON",
			"Set Up Threshold", "Set DownThreshold", "Configure SMTP Notification",
			"Configure Callback Notification", "Finish"}, 1, 10, false)

		switch selected {
		case 1:
			for {
				name := askNoSpaceNotEmpty("New Probe Name", "Name", stringDefault(probe.Name, helper.RandomName()), true)
				// todo check probe name duplicate
				probe.Name = name
				break
			}
		case 2:
			for {
				id := askNoSpaceNotEmpty("New Probe ID", "ID", stringDefault(probe.ID, uuid.New().String()), true)
				// todo check probe id duplicate
				probe.ID = id
				break
			}
		case 3:
			manageProbeRequest(probe)
		case 4:
			probe.BaseURL = interact.Ask("New Probe Base URL", stringDefault(probe.BaseURL, "http://localhost"), true)
		case 5:
			for {
				c := interact.Ask("New Probe Cron", stringDefault(probe.Cron, "0 */5 * * * * *"), true)
				_, err := cron.NewSchedule(c)
				if err != nil {
					fmt.Println("Cron syntax invalid")
					continue
				}
				// todo make some schedule rapid check
				probe.Cron = c
				break
			}
		case 6:
			probe.UpThreshold = interact.AskNumber("New Up Threshold", 2, 10, intDefault(probe.UpThreshold, 3), false)
		case 7:
			probe.DownThreshold = interact.AskNumber("New Down Threshold", 2, 10, intDefault(probe.DownThreshold, 3), false)
		case 8:
			configureSNMPNotification(probe)
		case 9:
			configureCallbackNotification(probe)
		case 10:
			return
		}
	}
}

func manageProbeRequest(probe *internal.Probe) {
	for {
		fmt.Printf("\n---[ PROBE \"%s\" REQUEST LIST ]-- (Base URL : %s)---------------\n", probe.Name, probe.BaseURL)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NO", "PROBE REQueST NAME", "PATH", "METHOD"})
		for idx, prob := range probe.Requests {
			table.Append([]string{fmt.Sprintf("%d", idx), prob.Name, prob.PathExpr, prob.MethodExpr})
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Add Request", "Remove Request", "View/Update Request", "Finish"}, 1, 4, false)

		requestList := make([]string, len(probe.Requests))
		for idx, pr := range probe.Requests {
			requestList[idx] = pr.Name
		}

		switch selected {
		case 1:
			addProbeRequest(probe)
		case 2:
			if len(requestList) == 0 {
				fmt.Println("There are no request to remove.")
			} else {
				requestList = append(requestList, "Cancel")
				sel := interact.Select("Select request to remove", requestList, 0, len(requestList)-1, false)
				if sel == len(requestList)-1 {
					continue
				}
				if interact.Confirm(fmt.Sprintf("Are you sure to remove probe %s", probe.Requests[sel].Name), false) {
					probe.Requests = append(probe.Requests[:sel], probe.Requests[sel+1:]...)
					requestList = append(requestList[:sel], requestList[sel+1:]...)
					continue
				}
			}
		case 3:
			if len(requestList) == 0 {
				fmt.Println("There are no request to view/edit.")
			} else {
				probeList := append(requestList, "Cancel")
				sel := interact.Select("Select request to edit", requestList, 0, len(requestList)-1, false)
				if sel == len(probeList)-1 {
					continue
				}
				editProbeRequest(probe, sel)
			}
		case 4:
			return
		}
	}
}

func addProbeRequest(probe *internal.Probe) {
	var newName, newPath, newMethod string
	for {
		newName = askNoSpaceNotEmpty("New request name?", "Name", helper.RandomName(), true)
		for _, p := range probe.Requests {
			if p.Name == newName {
				fmt.Println("That name already exist in list of requests")
				continue
			}
		}
		break
	}
	newPath = askNoSpaceNotEmpty("New request Path Expression?", "Path", fmt.Sprintf("\"/%s\"", helper.RandomName()), true)
	switch interact.Select("HTTP Method Expression ?", []string{"\"GET\"", "\"POST\"", "\"PUT\"", "\"DELETE\"", "\"PATCH\"", "\"HEAD\"", "\"OPTIONS\""}, 1, 1, false) {
	case 1:
		newMethod = "\"GET\""
	case 2:
		newMethod = "\"POST\""
	case 3:
		newMethod = "\"PUT\""
	case 4:
		newMethod = "\"DELETE\""
	case 5:
		newMethod = "\"PATCH\""
	case 6:
		newMethod = "\"HEAD\""
	case 7:
		newMethod = "\"OPTIONS\""
	}
	p := &internal.ProbeRequest{
		Name:        newName,
		PathExpr:    newPath,
		MethodExpr:  newMethod,
		HeadersExpr: make(map[string][]string),
	}
	probe.Requests = append(probe.Requests, p)
	editProbeRequest(probe, len(probe.Requests)-1)
}

func editProbeRequest(probe *internal.Probe, idx int) {
	probeRequest := probe.Requests[idx]

	for {
		fmt.Printf("\n---[ REQUEST %s CONFIGURATION ]--( %s )----------------\n", probeRequest.Name, probeRequest.PathExpr)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "VALUE"})

		if len(probeRequest.Name) > 0 {
			table.Append([]string{"Name", probeRequest.Name})
		} else {
			table.Append([]string{"Name", "Not Set"})
		}
		if len(probeRequest.PathExpr) > 0 {
			table.Append([]string{"Path Expression", probeRequest.PathExpr})
		} else {
			table.Append([]string{"Path Expression", "Not Set"})
		}
		if len(probeRequest.MethodExpr) > 0 {
			table.Append([]string{"Method Expression", probeRequest.MethodExpr})
		} else {
			table.Append([]string{"Method Expression", "Not Set"})
		}
		if len(probeRequest.BodyExpr) > 0 {
			table.Append([]string{"Body Expression", probeRequest.BodyExpr})
		} else {
			table.Append([]string{"Body Expression", "Not Set"})
		}

		table.Append([]string{"Request Headers", fmt.Sprintf("%d configured", len(probeRequest.HeadersExpr))})

		if len(probeRequest.StartRequestIfExpr) > 0 {
			table.Append([]string{"Start Request Criteria Expression", probeRequest.BodyExpr})
		} else {
			table.Append([]string{"Start Request Criteria  Expression", "Not Set"})
		}
		if len(probeRequest.SuccessIfExpr) > 0 {
			table.Append([]string{"Success Criteria Expression", probeRequest.SuccessIfExpr})
		} else {
			table.Append([]string{"Success Criteria Expression", "Not Set"})
		}
		if len(probeRequest.FailIfExpr) > 0 {
			table.Append([]string{"Fail Criteria Expression", probeRequest.FailIfExpr})
		} else {
			table.Append([]string{"Fail Criteria Expression", "Not Set"})
		}
		if len(probeRequest.CertificateCheckExpr) > 0 {
			table.Append([]string{"Certificate Check Expression", probeRequest.CertificateCheckExpr})
		} else {
			table.Append([]string{"Certificate Check Expression", "Not Set"})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Set Request Name",
			"Set Path Expression",
			"Set Method Expression",
			"Manage Header Expressions",
			"Set Body Expression",
			"Set Certificate Check Expression",
			"Set Start Request Criteria Expression",
			"Set Success Criteria Expression",
			"Set Fail Criteria Expression",
			"Finish",
		}, 1, 10, false)

		switch selected {
		case 1:
			for {
				name := askNoSpaceNotEmpty("New Request Name", "Name", stringDefault(probeRequest.Name, helper.RandomName()), true)
				// todo check request name duplicate
				probeRequest.Name = name
				break
			}
		case 2:
			for {
				pathExpr := interact.Ask("New Path Expression \n Must return a string", stringDefault(probeRequest.PathExpr, fmt.Sprintf("\"/%s\"", helper.RandomName())), true)
				if len(pathExpr) == 0 {
					fmt.Printf("Path expression can not be empty\n")
					continue
				}
				probeRequest.PathExpr = pathExpr
				break
			}
		case 3:
			switch interact.Select("HTTP Method Expression ?", []string{"\"GET\"", "\"POST\"", "\"PUT\"", "\"DELETE\"", "\"PATCH\"", "\"HEAD\"", "\"OPTIONS\""}, 1, 1, false) {
			case 1:
				probeRequest.MethodExpr = "\"GET\""
			case 2:
				probeRequest.MethodExpr = "\"POST\""
			case 3:
				probeRequest.MethodExpr = "\"PUT\""
			case 4:
				probeRequest.MethodExpr = "\"DELETE\""
			case 5:
				probeRequest.MethodExpr = "\"PATCH\""
			case 6:
				probeRequest.MethodExpr = "\"HEAD\""
			case 7:
				probeRequest.MethodExpr = "\"OPTIONS\""
			}
		case 4:
			manageRequestHeaders(probeRequest, probeRequest.PathExpr, probeRequest.MethodExpr)
		case 5:
			bodyExp := interact.Ask("New Body Expression? \n Must return string", stringDefault(probeRequest.BodyExpr, fmt.Sprintf("\"FooBar\"")), true)
			probeRequest.BodyExpr = bodyExp
		case 6:
			expr := interact.Ask("New Certificate Check Expression? \n Must return boolean", stringDefault(probeRequest.BodyExpr, fmt.Sprintf("\"FooBar\"")), true)
			probeRequest.CertificateCheckExpr = expr
		case 7:
			expr := interact.Ask("New Start Request Criteria Expression? \n Must return boolean", stringDefault(probeRequest.BodyExpr, fmt.Sprintf("\"FooBar\"")), true)
			probeRequest.StartRequestIfExpr = expr
		case 8:
			expr := interact.Ask("New Success Criteria Expression? \n Must return boolean", stringDefault(probeRequest.BodyExpr, fmt.Sprintf("\"FooBar\"")), true)
			probeRequest.SuccessIfExpr = expr
		case 9:
			expr := interact.Ask("New Fail Criteria Expression? \n Must return boolean", stringDefault(probeRequest.BodyExpr, fmt.Sprintf("\"FooBar\"")), true)
			probeRequest.FailIfExpr = expr
		case 10:
			return
		}
	}
}

func manageRequestHeaders(pr *internal.ProbeRequest, path, method string) {
	for {
		fmt.Printf("\n---[ REQUEST HEADERS ]--(%s)-(%s)------------------\n", path, method)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"KEY", "VALUES"})

		for k, v := range pr.HeadersExpr {
			table.Append([]string{k, strings.Join(v, ",")})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Add/Modify Header",
			"Remove Header",
			"Finish",
		}, 1, 3, false)
		switch selected {
		case 1:
			options := make([]string, 0)
			for k, _ := range pr.HeadersExpr {
				options = append(options, k)
			}
			options = append(options, "New Header", "Cancel")
			selected := interact.Select("Choose header to modify or select New Header or Cancel", options, 0, len(options)-1, false)
			switch selected {
			case len(options) - 1:
				// do nothing
			case len(options) - 2:
				key := askNoSpaceNotEmpty("Header Key ?", "Header Key", helper.RandomName(), true)
				values := interact.Ask("Header Value Expression (separated by semi-colon ';') ?\n Each expression must result a string", fmt.Sprintf("\"%s\"", helper.RandomName()), true)
				pr.HeadersExpr[key] = strings.Split(values, ";")
			default:
				key := options[selected]
				defaValues := pr.HeadersExpr[key]
				values := interact.Ask("Header Value Expression (separated by semi-colon ';') ?\n Each expression must result a string", strings.Join(defaValues, ";"), true)
				pr.HeadersExpr[key] = strings.Split(values, ";")
			}
		case 2:
			options := make([]string, 0)
			for k, _ := range pr.HeadersExpr {
				options = append(options, k)
			}
			options = append(options, "Cancel")
			selected := interact.Select("Choose header to remove or Cancel", options, 0, len(options)-1, false)
			switch selected {
			case len(options) - 1:
				// do nothing
			default:
				key := options[selected]
				delete(pr.HeadersExpr, key)
			}
		case 3:
			return
		}
	}
}

func configureSNMPNotification(probe *internal.Probe) {
	smtpConfig := probe.SMTPNotification
	if smtpConfig == nil {
		smtpConfig = &internal.SMTPNotificationTarget{
			To:  make([]*internal.Mailbox, 0),
			Cc:  make([]*internal.Mailbox, 0),
			Bcc: make([]*internal.Mailbox, 0),
		}
		probe.SMTPNotification = smtpConfig
	}
	for {
		fmt.Printf("\n---[ PROBE SMTP NOGIFICATION ]---------------------------\n")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ITEM", "VALUE"})

		if len(smtpConfig.SMTPHost) > 0 {
			table.Append([]string{"SMTP host", smtpConfig.SMTPHost})
		} else {
			table.Append([]string{"SMTP host", "Not Set"})
		}
		if smtpConfig.SMTPPort == 0 {
			table.Append([]string{"SMTP Port", fmt.Sprintf("%d", smtpConfig.SMTPPort)})
		} else {
			table.Append([]string{"SMTP Port", "Not Set"})
		}
		if smtpConfig.From != nil {
			table.Append([]string{"From", smtpConfig.From.String()})
		} else {
			table.Append([]string{"From", "Not Set"})
		}
		if len(smtpConfig.Password) > 0 {
			table.Append([]string{"Password", smtpConfig.Password})
		} else {
			table.Append([]string{"Password", "Not Set"})
		}
		if len(smtpConfig.To) > 0 {
			table.Append([]string{"To", fmt.Sprintf("%d mailboxes.", len(smtpConfig.To))})
		} else {
			table.Append([]string{"To", "No Recipient"})
		}
		if len(smtpConfig.Cc) > 0 {
			table.Append([]string{"To", fmt.Sprintf("%d mailboxes.", len(smtpConfig.Cc))})
		} else {
			table.Append([]string{"To", "No Recipient in CC"})
		}
		if len(smtpConfig.Bcc) > 0 {
			table.Append([]string{"To", fmt.Sprintf("%d mailboxes.", len(smtpConfig.Bcc))})
		} else {
			table.Append([]string{"To", "No Recipient in BCC"})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		selected := interact.Select("What to do ?", []string{
			"Set SMTP Host",
			"Set SMTP Port",
			"Set From",
			"Set Password",
			"Manage Recipients",
			"Manage Cc Recipients",
			"Manage Bcc Recipients",
			"Finish",
		}, 1, 8, false)
		switch selected {
		case 1:
			smtpConfig.SMTPHost = askNoSpaceNotEmpty("SMTP Host ?", "SMTP Host", stringDefault(smtpConfig.SMTPHost, "127.0.0.1"), true)
		case 2:
			smtpConfig.SMTPPort = interact.AskNumber("SMTP Port ?", 1, 65000, intDefault(smtpConfig.SMTPPort, 25), true)
		case 3:
			for {
				mailbox := interact.Ask("From Email ?", "Is Me<itsme@host.com>", true)
				mb, err := internal.NewMailbox(mailbox)
				if err != nil {
					fmt.Println("Invalid mail box. Should be in the format of 'mailbox@domain' or 'Display<mailbox@domain>'")
					fmt.Println("The email domain will be checked via MX lookup")
					continue
				}
				smtpConfig.From = mb
				break
			}
		case 4:
			smtpConfig.Password = interact.Ask("Password", "this is a very secret passphrase", true)
		case 5:
			manageEmailRecipients("TO", smtpConfig.To)
		case 6:
			manageEmailRecipients("CC", smtpConfig.Cc)
		case 7:
			manageEmailRecipients("BCC", smtpConfig.Bcc)
		case 8:
			return
		}
	}
}

func manageEmailRecipients(group string, mailBoxes []*internal.Mailbox) {
	for {
		fmt.Printf("\n---[ %s RECIPIENTS ]---------------------------\n", group)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"IDX", "NAME", "EMAIL"})
		for idx, mb := range mailBoxes {
			if mb.Name == "" {
				table.Append([]string{fmt.Sprintf("%d", idx), "noname", mb.Email})
			} else {
				table.Append([]string{fmt.Sprintf("%d", idx), mb.Name, mb.Email})
			}
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		options := []string{
			"Add Mailbox",
			"Remove Mailbox",
			"Finish",
		}
		switch interact.Select("What to do?", options, 1, 3, false) {
		case 1:
			for {
				mailbox := interact.Ask("Email ?", "Is Me<itsme@host.com>", true)
				mb, err := internal.NewMailbox(mailbox)
				if err != nil {
					fmt.Println("Invalid mail box. Should be in the format of 'mailbox@domain' or 'Display<mailbox@domain>'")
					fmt.Println("The email domain will be checked via MX lookup")
					continue
				}
				duplicate := false
				for _, mbox := range mailBoxes {
					if mbox.Email == mb.Email {
						fmt.Println("Email already a recipient, Updating Name")
						mbox.Name = mb.Name
						duplicate = true
						break
					}
				}
				if !duplicate {
					mailBoxes = append(mailBoxes, mb)
				}
				break
			}
		case 2:
			sbox := interact.Ask("Name or Email ?", "itsme@host.com", true)
			todel := -1
			for idx, mbox := range mailBoxes {
				if mbox.Email == sbox || mbox.Name == sbox {
					if interact.Confirm(fmt.Sprintf("Are you sure to delete %s", mbox.String()), false) {
						todel = idx
					}
					break
				}
			}
			mailBoxes = append(mailBoxes[:todel], mailBoxes[todel+1:]...)
		case 3:
			return
		}
	}
}

func configureCallbackNotification(p *internal.Probe) {
	if p.CallbackNotification == nil {
		p.CallbackNotification = &internal.CallbackNotificationTarget{}
	}
	for {
		fmt.Printf("\n---[ PROBE CALLBACK NOTIFICATION ]-----------------------\n")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "URL"})

		if len(p.CallbackNotification.UpCall) > 0 {
			table.Append([]string{"UP Callback URL", p.CallbackNotification.UpCall})
		} else {
			table.Append([]string{"UP Callback URL", "not specified"})
		}
		if len(p.CallbackNotification.DownCall) > 0 {
			table.Append([]string{"DOWN Callback URL", p.CallbackNotification.DownCall})
		} else {
			table.Append([]string{"DOWN Callback URL", "not specified"})
		}

		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()

		options := []string{
			"Set UP Callback URL",
			"Set DOWN Callback URL",
			"Finish",
		}

		switch interact.Select("What to do ?", options, 1, 3, false) {
		case 1:
			p.CallbackNotification.UpCall = interact.Ask("UP Callback URL?", stringDefault(p.CallbackNotification.UpCall, "http://localhost"), true)
		case 2:
			p.CallbackNotification.DownCall = interact.Ask("DOWN Callback URL?", stringDefault(p.CallbackNotification.DownCall, "http://localhost"), true)
		case 3:
			return
		}
	}
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
			central.ListenHost = interact.Ask("Specify New Server Host IP", stringDefault(central.ListenHost, "0.0.0.0"), true)
		case 2:
			central.ListenPort = interact.AskNumber("Pick New Port Number", 1025, 65000, intDefault(central.ListenPort, 53421), true)
		case 3:
			central.ReadHeaderTimeoutSecond = interact.AskNumber("Specify New Read-Header-Timeout in Second", 3, 120, intDefault(central.ReadHeaderTimeoutSecond, 10), true)
		case 4:
			central.ReadTimeoutSecond = interact.AskNumber("Specify New Read-Timeout in Second", 3, 120, intDefault(central.ReadTimeoutSecond, 10), true)
		case 5:
			central.WriteTimeoutSecond = interact.AskNumber("Specify New Write-Timeout in Second", 3, 120, intDefault(central.WriteTimeoutSecond, 10), true)
		case 6:
			central.IdleTimeoutSecond = interact.AskNumber("Specify New Idle-Timeout in Second", 3, 120, intDefault(central.IdleTimeoutSecond, 10), true)
		case 7:
			central.AdminUser = interact.Ask("Specify New Admin account name", stringDefault(central.AdminUser, "Admin"), true)
		case 8:
			central.AdminPassword = interact.Ask("Specify New Admin Password", stringDefault(central.AdminPassword, "This is a very good pass Phrase"), true)
		case 9:
			central.JWTIssuer = interact.Ask("Specify New JWT Issuer", stringDefault(central.JWTIssuer, "mihp.io"), true)
		case 10:
			central.JWTKey = interact.Ask("Specify New JWT Secret Key", stringDefault(central.JWTKey, "Th1sk3ymu$tb3ch4ngebef0r3use@ge0npr0duction"), true)
		case 11:
			central.JWTAccessKeyAgeMinute = interact.AskNumber("Specify New Access Token Age in Minute", 3, 60*24*365*10, intDefault(central.JWTAccessKeyAgeMinute, 10), true)
		case 12:
			central.JWTRefreshKeyAgeMinute = interact.AskNumber("Specify New Refresh Token Age in Minute", 3, 60*24*365*10, intDefault(central.JWTRefreshKeyAgeMinute, 60*24*365*2), true)
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

func stringDefault(tocheck, alternative string) string {
	if len(tocheck) == 0 {
		return alternative
	}
	return tocheck
}

func intDefault(tocheck, alternative int) int {
	if tocheck == 0 {
		return alternative
	}
	return tocheck
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
			cfg.Host = interact.Ask("Specify New DB Host IP", stringDefault(cfg.Host, "0.0.0.0"), true)
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
			cfg.User = interact.Ask("Specify New DB User", stringDefault(cfg.User, "root"), true)
		case 4:
			cfg.Password = interact.Ask("Specify New DB Password", stringDefault(cfg.Password, "this is db user password"), true)
		case 5:
			cfg.Database = interact.Ask("Specify New DB Schema", stringDefault(cfg.Database, "mihp"), true)
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
