# Valkyrie Launcher

The Launcher tool of Valkyrie is essentially a worker orchestration tool that recieves requests, and spins up worker containers to handle the task. 
The Launcher tool runs on any worker nodes (nodes that will launch containers, to handle work), and dispatches the container.
Realistically, this should never be called directly.  The Foreman should ferret out the Action Messages. Using Foreman for this provides you round-robin
dispatching of new worker containers, and duplication protection.  If you want to trigger these directly however, you can.

### /trigger
This takes an Action Message as well (see Foreman README.md) and actually spins up worker containers to handle the action message work.

### Environment Variables

#### loggerurl http://localhost:8092/
This is the URL/Port of the Valkyrie Logger Service.  All of the services in Valkyrie try to contact this URL to log it's output/actions

#### usesecurelogging false
Secure logging allows the logger to require a key value to allow things to write to it, this true/false value determines if you want to use that or not

#### logkey ""
If you are using secure logging, this is the key used to allow the logger to accept logging messages from this service

#### defaultimage worker
This is the default docker image to use if none in specified in the Action Message
