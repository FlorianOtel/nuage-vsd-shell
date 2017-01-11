# nuage-vsd-shell

A CLI shell for interacting with Nuage Networks VSD (Virtual Services Director). 

It is written as a wrapper around [Golang libraries for Nuage Networks VSD API](https://github.com/nuagenetworks/vspk-go/). Based on [abiosoft/ishell](https://github.com/abiosoft/ishell).

This is provided "as such", without any official support.

Thanks in advance for feedback, questions, raising issues, or contributing -- much appreciated.


## Usage

There are two types of shell commands:
* Auxiliary commands for E.g.: Displaying the details of (existing) API connection. Currently a single API connection is supported at one time; Setting the debug
 level (only two levels supported -- a verbose "Debug" level and "Info"); Setting the API connection details (API endpoint and credentials) and initializing a connection

* Wrappers around the Nuage Networks API calls themselves: "GET", "CREATE", "DELETE" etc. See below for commands currently supported.

### Auxiliary commands

```
Nuage API Interactive Shell
>> help
Commands:
CREATE DELETE GET clear debuglevel displayconn exit greet help makeconn setconn

>> debuglevel
Set debug level: Debug, Info (default): Debug
Debug level now set to: Debug

>> displayconn

    Endpoint URL: [https://127.0.0.1:8443]
    API version: [v4_0]
    Not connected

>> setconn
Set Nuage API connection Details: Endpoint IP address ; User + Password ; Nuage API version
  Enter your VSD IP address> 127.0.0.1
  Enter your Enterprise (organization) name. Leave empty if default > org
  Enter your username. Leave empty if default > user
  Enter your password. Leave empty if default > ****


### API wrapper commands

Each command has a 1-1 correspondence with the underlying library calls.

They require that a valid API connection is established first (using the `makeconn` command above).

Currently the following shell commands are supported. For more details and examples on how they use the API library please see `main.go`

```
#### GET operations


GET enterprises
GET enterprises <ID>
GET enterprises <ID> domaintemplates
GET enterprises <ID> domains
GET enterprises <ID> L2domains
GET enterprises <ID> vms
GET enterprises <ID> containers

GET domaintemplates <ID>
GET domaintemplates <ID> zonetemplates

GET domains
GET domains <ID>
GET domains <ID> vports
GET domains <ID> vminterfaces

GET zones
GET zones <ID>
GET subnets
GET subnets <ID>
GET subnets <ID> <child>

GET subnets <ID> vports
GET subnets <ID> vminterfaces

GET vms
GET vms <ID>

GET containers

#### DELETE operations

DELETE enterprise <ID>

DELETE domaintemplate <ID>

DELETE zonetemplate <ID>

DELETE domain <ID>

DELETE zone <ID>

DELETE subnet <ID>

DELETE vport <ID>

DELETE vminterface <ID>

DELETE vm <ID>

DELETE container <ID>
```


Example: Obtaining the list of organizations (enterprises) currently defined:

```
>> GET enterprises

 ===> Org nr [0]: Name [ORG1] <===
{
	"allowedForwardingClasses": [
				    "H"
				    ],
				    "associatedEnterpriseSecurityID": "429fb71e-9613-4ff6-80a6-dda00075b1e4",
				    "associatedGroupKeyEncryptionProfileID": "f915f3a7-8314-4a24-b853-51e913ba4e16",
				    "associatedKeyServerMonitorID": "f9a88422-548f-407e-aa00-146fcd479f7c",
				    "customerID": 10004,
				    "description": "20150911",
				    "DHCPLeaseinterval": 24,
				    "encryptionManagementMode": "DISABLED",
				    "enterpriseProfileID": "f1e5eb19-c67a-4651-90c1-3f84e23e1d36",
				    "floatingIPsQuota": 16,
				    "name": "ORG1",
				    "receiveMultiCastListID": "081169f6-cb2f-4c6e-8e94-b701224a5141",
				    "sendMultiCastListID": "738446cc-026f-488f-9718-b13f4390857b",
				    "creationDate": 1442002567000,
				    "lastUpdatedBy": "8a6f0e20-a4db-4878-ad84-9cc61756cd5e",
				    "lastUpdatedDate": 1442002905000,
				    "owner": "8a6f0e20-a4db-4878-ad84-9cc61756cd5e",
				    "entityScope": "ENTERPRISE",
				    "ID": "ea6862a3-b215-4343-b54f-3cee1d9ef9be"
}


 ===> Org nr [1]: Name [ORG2] <===
{
	"allowedForwardingClasses": [
				    "H"
				    ],
				    "associatedEnterpriseSecurityID": "626af6e0-74f8-4008-8122-ca89a2e40c28",
				    "associatedGroupKeyEncryptionProfileID": "a82f8626-772d-4505-a4e5-9b22c1daca63",
				    "associatedKeyServerMonitorID": "7788e5e5-91d1-4f9b-ab13-14efff7784bb",
				    "customerID": 10005,
				    "description": "Test org #2",
				    "DHCPLeaseinterval": 24,
				    "encryptionManagementMode": "DISABLED",
				    "enterpriseProfileID": "f1e5eb19-c67a-4651-90c1-3f84e23e1d36",
				    "floatingIPsQuota": 16,
				    "name": "ORG2",
				    "receiveMultiCastListID": "081169f6-cb2f-4c6e-8e94-b701224a5141",
				    "sendMultiCastListID": "738446cc-026f-488f-9718-b13f4390857b",
				    "creationDate": 1443106911000,
				    "lastUpdatedBy": "8a6f0e20-a4db-4878-ad84-9cc61756cd5e",
				    "lastUpdatedDate": 1443106911000,
				    "owner": "8a6f0e20-a4db-4878-ad84-9cc61756cd5e",
				    "entityScope": "ENTERPRISE",
				    "ID": "30d5ac86-edfa-49e0-afcc-786b63b41e9a"
}
Enterprise list -- done

>>
```
