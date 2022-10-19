# http-slave

This tool designed for remote command execution.

It periodically gets command execution request by HTTP from master server, execute the command and POST response to specified URL.

I used this tool to investigate a problem on remote host which periodically lose VPN connection and became inaccessible from the outside of NAT.

## How to use

Start http-slave on slave host:

    ./http-slave -i 60 -u https://your.domain/example.json

Provide example.json response somehow (nginx?) by master host:

    {
        "command": ["/usr/local/sbin/http-slave-chekscript.sh"],
        "respond_url": "https://youd.domain/response?foo=bar",
        "immediately_next": false
    }

Handle somehow POST response from http-slave somewhere.