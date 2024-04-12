# go-ports

Ports-go is a Network Monitor that shows real-time network connection tracking.
It lists active network connections along with associated processes, updating
this information periodically. The tool provides insights into network activity
on your system.

## Description

This tool scans and displays active network connections on your system. For each
connection, it shows details like protocol, local and remote addresses,
connection state, and the associated process. It's designed to update this
information in real-time, allowing users to monitor their network activity
continuously. The output is displayed in a terminal-based user interface and can
be exited gracefully with a Control+C command.

### Features

Real-time monitoring of network connections. Displays protocol, local and remote
IP addresses, state, and associated process name. Color-coded output for better
readability (on compatible terminals). Updates data periodically. Graceful exit
using Control+C. Installation To install and run Go Network Monitor, follow
these steps:

Ensure you have Go installed on your system. Download Go if you haven't
installed it yet.

Clone the repository or download the source code.

Navigate to the source code directory.

Compile the program with `go build` Run the executable.

On Unix-like systems: Run the executable. `./ports-go` On Windows:
`ports-go.exe`

## Usage

Once running, the program will display a table of active network connections
updating in real-time. Press `Control+C` at any time to exit the program.

### Credits

this is a small tribute to another similiar application called ports, which
pretty much a simple script in bash that once upon a time helped me figure out
why applications used what ports and why and it was pretty awesome.

