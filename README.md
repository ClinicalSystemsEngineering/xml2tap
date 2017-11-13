# xml2tap
r5 xml to tap output

*note this is a custom solution not necessarily supported by the vendor Rauland Borg
though the solution is very straightforward and simple it is currently undergoing extensive testing for integration into 
Connexall.
If you choose to use/test this application please provide feedback where applicable for future updates.

This application is a simple attemt for replacing the default Responder 5 paging service.
We have documented that the default Responder 5 paging service only supports a throughput of
about 38 messages/min.  In larger hospitals, traffic can easily overburden the system causing
messages to queue for several minutes.  This is especially prevalent during power surges that cause multiple bed exit
messages to generate or if staff terminal or other call assignments have multiple caregivers assigned
to a call stop.  As you can imagine this could be problematic for high priority notifications.

There are several command line flags that can be used to set various ports.
The following are the default settings:
-xmlPort 5051
-tapPort 10001
-httpPort 80

example useage: .\xml2tap -xmlPort 5050 -tapPort 1000 -httpPort 85


The application will listen on the xmlPort for Responder 5 xml messages and output these messages on the tap Port.
Note the tap port is a server connection and not a client.

When no message data is being received the application will produce a tap output reconnect routine as a means of keep alive.

The http site is used for status diagnositics and monitoring if messages beging to queue.  Currently this is hardcoded at 300 messages
but can easily be changed in the source before compiling.

The lumberjack loging package is used to log data to file and the default location is /var/log/xml2tap/xml2tap.log.


