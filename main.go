package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/abiosoft/ishell"

	log "github.com/Sirupsen/logrus"

	"github.com/FlorianOtel/go-bambou/bambou"
	"github.com/FlorianOtel/vspk-go/vspk"
	"github.com/howeyc/gopass"
)

var (
	// Nuage API connection defaults. We need to keep them as global vars since commands can be invoked in whatever order.

	root      *vspk.Me
	mysession *bambou.Session

	// vsdurl, org, user, passwd string
	// Temporary defaults
	vsdurl    = "https://172.16.254.7:7443"
	user      = "csproot"
	passwd    = "csproot"
	org       = "csp"
	certfname = "/root/certlogin1.pem"
	keyfname  = "/root/certlogin1-Key.pem"
)

////////
////////
////////

func resetconn(args ...string) (string, error) {
	if mysession != nil {
		mysession.Reset()
	}
	// mysession = nil
	root = nil
	return "", nil
}

// Establish Nuage API connection using user + password

func makeconn(args ...string) (string, error) {

	if user == "" || passwd == "" || vsdurl == "" {
		return "Invalid VSD url / user / password. Please set connection details using `setconn` command ", nil
	}

	mysession, root = vspk.NewSession(user, passwd, org, vsdurl)

	// fmt.Printf("===> My VSD URL is: %s\n", vsdurl)

	// fmt.Printf("===> My Bambou session is: %#v\n", *mysession)

	// mysession.SetInsecureSkipVerify(true)

	err := mysession.Start()

	if err != nil {
		resetconn()
		fmt.Printf("Nuage API connection failed: ")
		return "", err
	} else {
		return "Nuage VSD connection established", nil
	}
}

// Establish Nuage API connection -- using certificates

func makecertconn(args ...string) (string, error) {

	if vsdurl == "" {
		return "Invalid VSD url.Please set connection details using `setconn` command ", nil
	}

	if cert, err := tls.LoadX509KeyPair(certfname, keyfname); err != nil {

		fmt.Printf("Loading TLS certificate and private key failed: ")
		return "", err
	} else {
		mysession, root = vspk.NewX509Session(&cert, vsdurl)
	}

	// fmt.Printf("===> My VSD URL is: %s\n", vsdurl)

	// fmt.Printf("===> My Bambou session is: %#v\n", *mysession)

	// mysession.SetInsecureSkipVerify(true)

	if err := mysession.Start(); err != nil {
		resetconn()
		fmt.Printf("Nuage TLS API connection failed: ")
		return "", err
	} else {
		return "Nuage VSD TLS connection established", nil
	}
}

// Displays Nuage session details
func displayconn(args ...string) (string, error) {
	if root == nil {
		return "Not Connected", nil
	} else {
		fmt.Printf("Nuage VSD connection established as:\n    VSD URL: [%s]\n    User: [%s]\n    Organization: [%s]\n", mysession.URL, mysession.Username, mysession.Organization)
		return "", nil
	}

}

// Set Nuage API connection details in top level vars
func setconn(args ...string) (string, error) {
	var (
		err      error
		vsdip    string
		bytepass []byte
	)

	fmt.Println("Set Nuage API connection Details: VSD IP address ; User ;  Password ")

	// Get VSD IP
	fmt.Print("\tEnter your VSD IP address > ")
	_, err = fmt.Scanln(&vsdip)

	if err == nil {
		if net.ParseIP(vsdip) == nil {
			return "'" + vsdip + "'" + " is not a valid IP address", nil
		}
		// We assume (hardcode) that the URL for the Nuage API has the form "https://<VSD_ip_addr>:8443"
		vsdurl = "https://" + vsdip + ":8443"
	} else {
		if err.Error() != "unexpected newline" { // If we just press Enter the VSD IP / URL remains the same as previous value
			return "Error: ", err
		}
	}

	// Get username
	fmt.Print("\tEnter your username > ")
	_, err = fmt.Scanln(&user)

	if err != nil && err.Error() != "unexpected newline" {
		return "Error: ", err
	}

	// Get password -- use "gopass"
	fmt.Print("\tEnter your password > ")
	bytepass, err = gopass.GetPasswdMasked()

	if err != nil && err.Error() != "unexpected newline" {
		return "Error: ", err
	}

	passwd = string(bytepass)

	// Get Enterprise name.
	fmt.Print("\tEnter your Enterprise (organization) name > ")
	_, err = fmt.Scanln(&org)

	if err != nil && err.Error() != "unexpected newline" {
		return "Error: ", err
	}

	// TBD: Insert code for changing the Nuage API version here. Currently only 4_0 (hardcoded)

	return "", nil

}

// Test function -- dummy greet
func mygreet(args ...string) (string, error) {
	name := "Stranger"
	if len(args) > 0 {
		name = strings.Join(args, " ")
	}
	return "Hello " + name, nil
}

// Set debug level
func debuglevel(args ...string) (string, error) {
	var loglevel string
	fmt.Print("Set debug level: Debug (debug), Info: ")
	_, err := fmt.Scanln(&loglevel)
	if err != nil {
		if err.Error() == "unexpected newline" {
			// log.SetLevel(log.InfoLevel)
			// return "Debug level now set to Info", nil
			log.SetLevel(log.DebugLevel)
			return "Debug level now set to Debug", nil
		}
		return "Error", err
	}

	switch loglevel {
	case "Debug":
		log.SetLevel(log.DebugLevel)
		return "Debug level now set to: " + loglevel, nil
	case "Info":
		log.SetLevel(log.InfoLevel)
		return "Debug level now set to: " + loglevel, nil
	default:
		log.SetLevel(log.InfoLevel)
		return "Debug level now set to Info", nil
	}
}

func main() {

	// First we set the connection details

	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.

	shell := ishell.New()

	shell.Println("Nuage VSD API Interactive Shell")

	shell.Register("greet", mygreet)

	shell.Register("debuglevel", debuglevel)

	// API connection handling

	shell.Register("setconn", setconn)

	shell.Register("makeconn", makeconn)

	shell.Register("makecertconn", makecertconn)

	shell.Register("displayconn", displayconn)

	shell.Register("resetconn", resetconn)

	//// Top-level CRUD operations
	shell.Register("GET", Get)
	//
	// shell.Register("CREATE", Create)
	//
	shell.Register("DELETE", Delete)

	// start shell
	shell.Start()
}

func Get(args ...string) (string, error) {

	if root == nil {
		return "Not Connected to a VSD server", nil
	}

	// 1 argument:  <entity>
	// 2 arguments: <entity> <ID>
	// 3 arguments: <entity> <ID> <children>

	if len(args) < 1 || len(args) > 3 {
		return "GET <entity> [ <ID> [ <children> ] ]  ", nil
	}

	entity := args[0]

	switch entity {
	case "containers": // GET containers
		containerlist, err := root.Containers(&bambou.FetchingInfo{})

		if err != nil {
			fmt.Printf("GET containers failed: ")
			return "", err
		}

		for i, v := range containerlist {
			container, _ := json.MarshalIndent(*v, "", "\t")
			fmt.Printf("\n ===> Container nr [%d]: Name [%s] <=== \n%#s\n", i, containerlist[i].Name, string(container))
		}

		return "Enterprise list -- done", nil

	case "enterprises":
		switch len(args) {
		case 1: // GET enterprises

			orglist, err := root.Enterprises(&bambou.FetchingInfo{})

			if err != nil {
				fmt.Printf("GET enterprises failed: ")
				return "", err
			}

			for i, v := range orglist {
				org, _ := json.MarshalIndent(*v, "", "\t")
				fmt.Printf("\n ===> Org nr [%d]: Name [%s] <=== \n%#s\n", i, orglist[i].Name, string(org))
			}

			return "Enterprise list -- done", nil

		case 2: // GET enterprises <ID>

			org := new(vspk.Enterprise)
			org.ID = args[1]
			err := org.Fetch()

			if err != nil {
				fmt.Printf("GET enterprise ID [%s] failed: ", org.ID)
				return "", err
			}

			// JSON pretty-print the org
			jsonorg, _ := json.MarshalIndent(org, "", "\t")
			fmt.Printf("\n\n ===> Org: Name [%s] <=== \n%#s\n", org.Name, string(jsonorg))

			return "Enterprise Get ID -- done", nil

		case 3: // GET enterprises <ID> <child>
			org := new(vspk.Enterprise)
			org.ID = args[1]
			child := args[2]

			switch child {
			case "domaintemplates": // GET enterprises <ID> domaintemplates
				// Get list of domain templates for that org

				dtl, err := org.DomainTemplates(&bambou.FetchingInfo{})

				if err != nil {
					fmt.Printf("GET enterprise [%s] domaintemplates failed: ", org.ID)
					return "", err
				}
				// Iterate through the list of domain templates and JSON pretty-print them
				fmt.Printf("\n ######## Domain templates for Enterprise ID: [%s] ########\n", org.ID)
				for i, v := range dtl {
					dt, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> Domain template nr [%d]: Name [%s] <=== \n%#s\n", i, dtl[i].Name, string(dt))
				}

				return "Domain template list -- done", nil

			case "domains": // GET enterprises <ID> domains
				// Get list of domains for the Enterprise ID

				dl, err := org.Domains(&bambou.FetchingInfo{})

				if err != nil {
					fmt.Printf("GET enterprise [%s] domains failed: ", org.ID)
					return "", err
				}

				// Iterate through the list of domains and JSON pretty-print them
				fmt.Printf("\n ######## Domains for Enterprise ID: [%s] ########\n", org.ID)
				for i, v := range dl {
					domain, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> Domain nr [%d]: Name [%s] <=== \n%#s\n", i, dl[i].Name, string(domain))
				}
				return "Domain list -- done", nil

			case "L2domains": // GET enterprises <ID> L2domains
				// Get list of domains for the Enterprise ID

				dl, err := org.L2Domains(&bambou.FetchingInfo{})

				if err != nil {
					fmt.Printf("GET enterprise [%s] L2domains failed: ", org.ID)
					return "", err
				}

				// Iterate through the list of domains and JSON pretty-print them
				fmt.Printf("\n ######## Domains for Enterprise ID: [%s] ########\n", org.ID)
				for i, v := range dl {
					domain, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> Domain nr [%d]: Name [%s] <=== \n%#s\n", i, dl[i].Name, string(domain))
				}
				return "Domain list -- done", nil

			case "vms": // GET enterprises <ID> vms
				// Get list of vms for the Enterprise ID

				dl, err := org.VMs(&bambou.FetchingInfo{})

				if err != nil {
					fmt.Printf("GET enterprise [%s] vms failed: ", org.ID)
					return "", err
				}

				// Iterate through the list of vms and JSON pretty-print them
				fmt.Printf("\n ######## Vms for Enterprise ID: [%s] ########\n", org.ID)
				for i, v := range dl {
					domain, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> Domain nr [%d]: Name [%s] <=== \n%#s\n", i, dl[i].Name, string(domain))
				}
				return "VM list -- done", nil

			case "containers": // GET enterprises <ID> containers
				// Get list of containers for the Enterprise ID

				dl, err := org.Containers(&bambou.FetchingInfo{})

				if err != nil {
					fmt.Printf("GET enterprise [%s] containers failed: ", org.ID)
					return "", err
				}

				// Iterate through the list of containers and JSON pretty-print them
				fmt.Printf("\n ######## Containers for Enterprise ID: [%s] ########\n", org.ID)
				for i, v := range dl {
					domain, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> Domain nr [%d]: Name [%s] <=== \n%#s\n", i, dl[i].Name, string(domain))
				}
				return "Container list -- done", nil
			}
		}

	case "domaintemplates":
		switch len(args) {
		case 2: // GET domaintemplates <ID>
			dt := new(vspk.DomainTemplate)
			dt.ID = args[1]
			err := dt.Fetch()

			if err != nil {
				fmt.Printf("GET domaintemplates ID [%s] failed. Error: ", dt.ID)
				return "", err
			}

			// JSON pretty-print the domain template
			jsondt, _ := json.MarshalIndent(dt, "", "\t")
			fmt.Printf("\n ===> Domain Template: Name [%s] <=== \n%#s\n", dt.Name, string(jsondt))
			return "Domain Template Get -- done", nil

		case 3: // GET domaintemplates <ID> <child>
			dt := new(vspk.DomainTemplate)
			dt.ID = args[1]
			child := args[2]

			switch child {
			case "zonetemplates": // GET domaintemplates <ID> zonetemplates

				ztl, err := dt.ZoneTemplates(&bambou.FetchingInfo{})

				if err != nil {
					fmt.Printf("GET domaintemplates ID [%s] zonetemplates failed. Error: ", dt.ID)
					return "", err
				}

				// Iterate through the list of zone templates and JSON pretty-print them
				fmt.Printf("\n ######## Zone templates for Domain template ID: [%s] ########\n", dt.ID)
				for i, v := range ztl {
					zt, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> Zone template nr [%d]: Name [%s] <=== \n%#s\n", i, ztl[i].Name, string(zt))
				}

				return "Zone template list -- done", nil
			}
		}

	case "domains":
		switch len(args) {
		case 1: // GET domains
			dl, err := root.Domains(&bambou.FetchingInfo{})

			if err != nil {
				fmt.Printf("GET domains failed: ")
				return "", err
			}

			for i, v := range dl {
				jsondomain, _ := json.MarshalIndent(v, "", "\t")
				fmt.Printf("\n ===> Domain nr [%d]: Name [%s] <=== \n%#s\n", i, dl[i].Name, string(jsondomain))
			}
			return "Domain list -- done", nil
		case 2: // GET domains <ID>
			domain := new(vspk.Domain)
			domain.ID = args[1]
			err := domain.Fetch()

			if err != nil {
				return "", err
			}

			jsondomain, _ := json.MarshalIndent(domain, "", "\t")
			fmt.Printf("\n ===> Domain Name [%s] <=== \n%#s\n", domain.Name, string(jsondomain))
			return "Domain Get -- done", nil

		case 3: // GET domains <ID> <child>
			domain := new(vspk.Domain)
			domain.ID = args[1]
			child := args[2]

			switch child {
			case "vports": // GET domains <ID> vports
				vports, err := domain.VPorts(&bambou.FetchingInfo{})

				if err != nil {
					return "", err
				}

				for i, v := range vports {
					jsonvport, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> VPort nr [%d]: Name [%s] <=== \n%#s\n", i, vports[i].Name, string(jsonvport))
				}
				return "Domain VPorts list -- done", nil

			case "vminterfaces": // GET domains <ID> vminterfaces
				vmiflist, err := domain.VMInterfaces(&bambou.FetchingInfo{})

				if err != nil {
					return "", err
				}

				for i, v := range vmiflist {
					jsonvmif, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> VMInterface nr [%d]: Name [%s] <=== \n%#s\n", i, vmiflist[i].Name, string(jsonvmif))
				}
				return "Subnet VMInterfaces list -- done", nil
			}
		}

	case "zones":
		switch len(args) {
		case 1: // GET zones
			zl, err := root.Zones(&bambou.FetchingInfo{})

			if err != nil {
				fmt.Printf("GET zones failed: ")
				return "", err
			}

			for i, v := range zl {
				jsonzone, _ := json.MarshalIndent(v, "", "\t")
				fmt.Printf("\n ===> Zone nr [%d]: Name [%s] <=== \n%#s\n", i, zl[i].Name, string(jsonzone))
			}

			return "Zone list -- done", nil

		case 2: // GET zones <ID>
			// Get a specific Zone ID

			zone := new(vspk.Zone)
			zone.ID = args[1]
			err := zone.Fetch()

			if err != nil {
				return "", err
			}

			jsonzone, _ := json.MarshalIndent(zone, "", "\t")
			fmt.Printf("\n ===> Zone Name [%s] <=== \n%#s\n", zone.Name, string(jsonzone))
			return "Zone Get -- done", nil
		}

	case "subnets":
		switch len(args) {
		case 1: // GET subnets
			// Get list of subnets with "nil" as parent domain -- i.e. global list of all subnets

			subnetlist, err := root.Subnets(&bambou.FetchingInfo{})

			if err != nil {
				fmt.Printf("GET subnets failed: ")
				return "", err
			}

			for i, v := range subnetlist {
				jsonsubnet, _ := json.MarshalIndent(v, "", "\t")
				fmt.Printf("\n ===> Subnet nr [%d]: Name [%s] <=== \n%#s\n", i, subnetlist[i].Name, string(jsonsubnet))
			}

			return "Subnet list -- done", nil

		case 2: // GET subnets <ID>
			// Get a specific Subnet ID
			subnet := new(vspk.Subnet)
			subnet.ID = args[1]
			err := subnet.Fetch()

			if err != nil {
				return "", err
			}

			jsonsubnet, _ := json.MarshalIndent(subnet, "", "\t")
			fmt.Printf("\n ===> Subnet Name [%s] <=== \n%#s\n", subnet.Name, string(jsonsubnet))
			return "Subnet Get -- done", err

		case 3: // GET subnets <ID> <child>
			subnet := new(vspk.Subnet)
			subnet.ID = args[1]
			child := args[2]

			switch child {
			case "vports": // GET subnets <ID> vports
				vports, err := subnet.VPorts(&bambou.FetchingInfo{})
				if err != nil {
					return "", err
				}

				for i, v := range vports {
					jsonvport, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> VPort nr [%d]: Name [%s] <=== \n%#s\n", i, vports[i].Name, string(jsonvport))
				}

				return "Subnet VPorts list -- done", nil

			case "vminterfaces": // GET subnets <ID> vminterfaces
				vmiflist, err := subnet.VMInterfaces(&bambou.FetchingInfo{})

				if err != nil {
					return "", err
				}

				for i, v := range vmiflist {
					jsonvmi, _ := json.MarshalIndent(v, "", "\t")
					fmt.Printf("\n ===> VMInterface nr [%d]: Name [%s] <=== \n%#s\n", i, vmiflist[i].Name, string(jsonvmi))
				}
				return "Subnet VMInterfaces list -- done", nil
			}
		}

	case "vms":
		switch len(args) {
		case 1: // GET vms
			vmlist, err := root.VMs(&bambou.FetchingInfo{})

			if err != nil {
				fmt.Printf("GET vms failed: ")
				return "", err
			}

			for i, v := range vmlist {
				jsonvm, _ := json.MarshalIndent(v, "", "\t")
				fmt.Printf("\n ===> VirtualMachine nr [%d]: Name [%s] <=== \n%#s\n", i, vmlist[i].Name, string(jsonvm))
			}

			return "VirtualMachine list -- done", nil

		case 2: // GET vms <ID>
			vm := new(vspk.VM)
			vm.ID = args[1]
			err := vm.Fetch()

			if err != nil {
				return "", err
			}

			jsonvm, _ := json.MarshalIndent(vm, "", "\t")
			fmt.Printf("\n ===> VirtualMachine Name [%s] <=== \n%#s\n", vm.Name, string(jsonvm))
			return "Virtual Machine Get -- done", nil
		}

	default:
		// Unknown entity request
		break
	}
	return "Don't know how to process Nuage API entity: " + strings.Join(args, " "), nil
}

func Delete(args ...string) (string, error) {
	if root == nil {
		return "Not Connected", nil
	}
	// Format: <entity> <ID>
	if len(args) != 2 {
		return "Format:\n    DELETE <entity> <ID>", nil
	}
	entity := args[0]
	id := args[1]

	switch entity {
	case "enterprise": // DELETE enterprise <ID>
		obj := new(vspk.Enterprise)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE enterprise ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "domaintemplate": // DELETE domaintemplate <ID>
		obj := new(vspk.DomainTemplate)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE domaintemplate ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "domain": // DELETE domain <ID>
		obj := new(vspk.Domain)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE domain ID [%s] failed. Error: ", obj.ID)
			return "", err
		}

	case "zonetemplate": // DELETE zonetemplate <ID>
		obj := new(vspk.ZoneTemplate)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE zonetemplate ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "zone": // DELETE zone <ID>
		obj := new(vspk.Zone)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE zone ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "subnet": // DELETE subnet <ID>
		obj := new(vspk.Subnet)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE subnet ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "vport": // DELETE vport <ID>
		obj := new(vspk.VPort)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE vport ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "vminterface": // DELETE vminterface <ID>
		obj := new(vspk.VMInterface)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE vminterface ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "vm": // DELETE vm <ID>
		obj := new(vspk.VM)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE vm ID [%s] failed. Error: ", obj.ID)
			return "", err
		}
	case "container": // DELETE container <ID>
		obj := new(vspk.Container)
		obj.ID = id
		err := obj.Delete()
		if err != nil {
			fmt.Printf("DELETE container ID [%s] failed. Error: ", obj.ID)
			return "", err
		}

	default:
		return "Don't know how to DELETE entity: " + entity, nil
	}
	return "", nil
}

//
// func Create(args ...string) (string, error) {
//
// 	// At least 2 arguments: entity <Name>
//
// 	if len(args) < 2 {
// 		return "Format:\n    CREATE <entity> <Name> [options]", nil
// 	}
//
// 	entity := args[0]
//
// 	switch entity {
// 	case "enterprise":
// 		if len(args) != 2 {
// 			return "Format:\n    CREATE enterprise <Name>", nil
// 		}
//
// 		// CREATE enterprise <Name>
// 		org := new(nuage_v3_2.Enterprise)
// 		org.Name = args[1]
// 		err := org.Create(myconn)
// 		if err != nil {
// 			return "", err
// 		}
//
// 		// JSON pretty-print the org
// 		jsonorg, _ := json.MarshalIndent(org, "", "\t")
// 		fmt.Printf("\n ===> Org: [%s] <=== \n%#s\n", org.Name, string(jsonorg))
// 		return "", err
//
// 	case "domaintemplate":
// 		if len(args) != 3 {
// 			return "Format:\n    CREATE domaintemplate <Name> <Parent Enterprise ID> ", nil
// 		}
//
// 		// CREATE domaintemplate <Name> <Parent Enterprise ID>
// 		dt := new(nuage_v3_2.Domaintemplate)
// 		dt.Name = args[1]
// 		dt.ParentID = args[2]
// 		err := dt.Create(myconn)
// 		if err != nil {
// 			return "", err
// 		}
// 		// JSON pretty-print the domain template
// 		jsondt, _ := json.MarshalIndent(dt, "", "\t")
// 		fmt.Printf("\n ===> Domain Template: Name [%s] <=== \n%#s\n", dt.Name, string(jsondt))
// 		return "Domain Template Create -- done", err
//
// 	case "domain":
// 		if len(args) != 4 {
// 			return "Format:\n    CREATE domain <Name> <Parent Enterprise ID> <Domain template ID>", nil
// 		}
// 		// CREATE domain <Name> <Parent Enterprise ID> <Domain template ID>
// 		domain := new(nuage_v3_2.Domain)
// 		domain.Name = args[1]
// 		domain.ParentID = args[2]
// 		domain.TemplateID = args[3]
// 		err := domain.Create(myconn)
// 		if err != nil {
// 			return "", err
// 		}
// 		jsondomain, _ := json.MarshalIndent(domain, "", "\t")
// 		fmt.Printf("\n ===> Domain Name [%s] <=== \n%#s\n", domain.Name, string(jsondomain))
// 		return "Domain Create -- done", err
//
// 	case "zonetemplate":
// 		if len(args) != 3 {
// 			return "Format:\n    CREATE zonetemplate <Name> <Parent domain template ID>", nil
// 		}
// 		// CREATE zonetemplate <Name> <Parent domain template ID>
// 		zt := new(nuage_v3_2.Zonetemplate)
// 		zt.Name = args[1]
// 		zt.ParentID = args[2]
// 		err := zt.Create(myconn)
// 		if err != nil {
// 			return "", err
// 		}
// 		jsonzt, _ := json.MarshalIndent(zt, "", "\t")
// 		fmt.Printf("\n ===> Zone template: Name [%s] <=== \n%#s\n", zt.Name, string(jsonzt))
// 		return "Zone Template Create -- done", err
// 	case "zone":
// 		if len(args) < 3 {
// 			return "Format:\n    CREATE zone <Name> <Parent Domain ID> [ <Zone template ID> ]", nil
// 		}
// 		// CREATE zone <Name> <Parent Domain ID> [ <Zone template ID> ]
// 		zone := new(nuage_v3_2.Zone)
// 		zone.Name = args[1]
// 		zone.ParentID = args[2]
// 		if len(args) >= 4 {
// 			zone.TemplateID = args[3]
// 		}
// 		err := zone.Create(myconn)
// 		if err != nil {
// 			return "", err
// 		}
// 		jsonzone, _ := json.MarshalIndent(zone, "", "\t")
// 		fmt.Printf("\n ===> Zone Name [%s] <=== \n%#s\n", zone.Name, string(jsonzone))
// 		return "Zone Create -- done", err
//
// 	case "subnet":
// 		if len(args) < 4 {
// 			return "Format:\n    CREATE subnet <Name> <Parent Zone ID> <Subnet template ID> \n or:\n     CREATE subnet <Name> <Parent Zone ID> <Subnet address> <Subnet mask>\n", nil
// 		}
// 		switch len(args) {
// 		case 4:
// 			// CREATE subnet <Name> <Parent Subnet ID> <Subnet template ID>
// 			subnet := new(nuage_v3_2.Subnet)
// 			subnet.Name = args[1]
// 			subnet.ParentID = args[2]
// 			subnet.TemplateID = args[3]
// 			err := subnet.Create(myconn)
// 			if err != nil {
// 				return "", err
// 			}
// 			jsonsubnet, _ := json.MarshalIndent(subnet, "", "\t")
// 			fmt.Printf("\n ===> Subnet Name [%s] <=== \n%#s\n", subnet.Name, string(jsonsubnet))
// 			return "Subnet Create -- done", err
// 		case 5:
// 			// CREATE subnet <Name> <Parent Subnet ID> <Subnet address> <Subnet mask>
// 			subnet := new(nuage_v3_2.Subnet)
// 			subnet.Name = args[1]
// 			subnet.ParentID = args[2]
// 			// TBD -- make sure these are proper dot notation...
// 			subnet.Address = args[3]
// 			subnet.Netmask = args[4]
// 			err := subnet.Create(myconn)
// 			if err != nil {
// 				return "", err
// 			}
// 			jsonsubnet, _ := json.MarshalIndent(subnet, "", "\t")
// 			fmt.Printf("\n ===> Subnet Name [%s] <=== \n%#s\n", subnet.Name, string(jsonsubnet))
// 			return "Subnet Create -- done", err
// 		}
//
// 	case "vport":
// 		if len(args) < 3 {
// 			return "Format:\n    CREATE vport <Name> <Parent Subnet ID> [ options ...]", nil
// 		}
// 		// CREATE vport <Name> <Parent Subnet ID> [ options ...]
// 		subnet := new(nuage_v3_2.Subnet)
// 		subnet.ID = args[2]
//
// 		var vport nuage_v3_2.VPort
// 		vport.Name = args[1]
// 		// ???
// 		vport.ID = vport.Name
// 		vport.Type = "VM"
// 		vport.AddressSpoofing = "INHERITED"
//
// 		vport.Active = true
// 		// ???? Not needed but still....
// 		vport.ParentID = subnet.ID
// 		vport.ParentType = "subnet"
//
// 		jsonvport, err := json.MarshalIndent(vport, "", "\t")
// 		fmt.Printf("\n ===> Created VPort: Name [%s] <=== \n%#s\n", vport.Name, string(jsonvport))
//
// 		vp, err := subnet.AddVPort(myconn, vport)
//
// 		if err != nil {
// 			return "", err
// 		}
//
// 		jsonvp, err := json.MarshalIndent(vp, "", "\t")
// 		fmt.Printf("\n ===> Created VPort: Name [%s] <=== \n%#s\n", vp.Name, string(jsonvp))
//
// 		return "VPort Create -- done", err
//
// 	case "vm":
// 		if len(args) != 5 {
// 			return "Format:\n    CREATE vm <Name> <UUID> <Interface0-MAC> <Interface0-VPortID>", nil
// 		}
// 		// CREATE vm <Name> <UUID> <Interface0-MAC> <Interface0-VPortID>
// 		var vm nuage_v3_2.VirtualMachine
// 		vm.Name = args[1]
// 		vm.UUID = args[2]
//
// 		var vmi nuage_v3_2.VMInterface
// 		vmi.MAC = args[3]
// 		vmi.VPortID = args[4]
//
// 		vm.Interfaces = append(vm.Interfaces, vmi)
//
// 		err := (&vm).Create(myconn)
// 		if err != nil {
// 			return "", err
// 		}
// 		jsonvm, _ := json.MarshalIndent(vm, "", "\t")
// 		fmt.Printf("\n ===> Virtual Machine: Name [%s] <=== \n%#s\n", vm.Name, string(jsonvm))
// 		return "Virtual Machine Create -- done", err
//
// 	default:
// 		return "Don't know how to create Nuage API entity: [" + entity + "]" + " with Name [" + args[1] + "]" + " and options: " + strings.Join(args[2:], " "), nil
// 	}
// 	return "", nil
// }
