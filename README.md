# nuage-vsd-shell

A CLI shell for interacting with Nuage Networks VSD (Virtual Services Director). 

It is written as a wrapper around [Golang libraries for Nuage Networks VSD API](https://github.com/nuagenetworks/vspk-go/). Based on [abiosoft/ishell](https://github.com/abiosoft/ishell)

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
  Enter your password. Leave empty if default > pass


