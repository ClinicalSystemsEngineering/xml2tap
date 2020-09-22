# xml2tap
r5 xml to tap output

*note this is a custom solution not necessarily supported by the vendor Rauland Borg
though the solution is very straightforward and simple it is currently undergoing extensive testing for integration into 
Connexall.
If you choose to use/test this application please provide feedback where applicable for future updates.

This application is a simple solution for replacing the default Responder 5 paging service.
We have documented that the default Responder 5 paging service only supports a throughput of
about 38 messages/min.  In larger hospitals, traffic can easily overburden the system causing
messages to queue for several minutes.  This is especially prevalent during power surges that can cause multiple bed exit
messages to generate or if staff terminal or other call assignments have multiple caregivers assigned
to a call stop.  As you can imagine this could be problematic for high priority notifications such as Code Blue.

There are several command line flags that can be used to set various ports.
The following are the default settings:


Flag: -xmlPort DefaultValue: 5051 Description: xml listener port for localhost

Flag: -tapPort DefaultValue: 10001 Description: localhost listener port for TAP server

Flag: -httpPort DefaultValue: 80 Description: localhost listner port for http server

Flag: -pprofPort DefaultValue: 8080 Description: localhost listner port for http server

Flag: -tapAdr DefaultValue: 127.0.0.1:10001 Description: server address for TAP client form is serverip:port

Flag: -tapwhitelist DefaultValue: 127.0.0.1 Description: ip address for incoming tap connection

Flag: -xmlwhitelist DefaultValue:127.0.0.1 Description:ip address for incoming xml connection")

example useage: .\xml2tap -xmlPort 5050 -tapPort 1000 -httpPort 85


The application will listen on the xmlPort for Responder 5 xml messages and output these messages on the tap Port.
Note the tap port is a server connection and not a client.

When no message data is being received the application will produce a tap output reconnect routine as a means of keep alive.

The http site is used for status diagnositics and monitoring if messages begin to queue.  Currently this is hardcoded at 10000 messages
but can easily be changed in the source before compiling.  This is found on the line

var parsedmsgs = make(chan string, 10000) //message processing channel for xml2tap conversions

Changing this for future use to be a command line flag is a trivial case, but was not implemented.

The lumberjack loging package is used to log data to file and the default location is /var/log/xml2tap/xml2tap.log.  This is also not a command line flag, but is left to be implemented as a trivial case.

The executable file was compiled on a windows 10 machine and has been tested on windows server 2012 r2 and windows 2016 with R5 version T14 and T15.

Recommend utilizing a service managment app to install the exe as a service i.e nssm (https://nssm.cc/)


