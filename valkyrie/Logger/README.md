# Valkyrie Logger 

The Logger tool of Valkyrie is a system to seperate out logging functionality for all the tools into a consolidated area.  The logger can be configured, and then all the other tools
simply call the logger to log messages.  

### /log
The /log end point is used to actually log the message from a Valkyrie Service.  It requires a few GET parameters:

#### Service
The service that is calling the logger. Such as "Alerter", or "Foreman"

#### LogLevel
The LogLevel (info, warning, error) this isn't really used.  It's here for future growth

#### Message
The Message to actually Log (action taken, failures, etc)

#### Key
The key passed in from a service if usesecurelogging is true/enabled


### Environment Variables

#### useeventstreams false
useeventstreams is an option that instead of logging to a file, it will throw all of its logdata to an open port on a server.  Splunk allows this as TCP or UDP Data Inputs.
This is restricted to TCP at the moment

#### eshostport splunk.yourdomain.com:5020
eshostport specifies the server and port number to connect to (via tcp) to send the bytes to

#### usesecurelogging false
This will enable secure logging so that clients/services require a key to log messages

#### logkey ""
This is the key that needs to be used on the clients/services in order to log to the service if it has secure logging enabled

### logfilelocation /tmp/valkyrie.log
This is the location to write the log file to if you aren't using event streams (OR if the event stream fails to write)

