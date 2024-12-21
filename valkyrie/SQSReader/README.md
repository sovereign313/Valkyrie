# Valkyrie SQSReader

The Valkyrie SQSReader is the cooperative sibling to the _dispatcher_.  The dispatcher takes incoming Action Messages (from URL Calls/Webhooks) and then (encrypts)
and puts messages into the AWS Queue.  The SQSReader is sort of like the Foreman, but gets the Action Messages from the AWS Queue instead of a direct HTTP Call.
This is really only necessary if you want to use AWS, or if you want to span the workers across datacenters, without opening ports for HTTP calls.  There are a few
benefits to using Dispatcher and SQSReader (as opposed to just the Foreman), such as scalability across data centers, but many people / companies are restricted to
what they can reach from behind their corporate Proxy/Firewall.

### Environment Variables

#### loggerurl http://localhost:8092/
This is the URL/Port of the Valkyrie Logger Service.  All of the services in Valkyrie try to contact this URL to log it's output/actions

#### usesecurelogging false
Secure logging allows the logger to require a key value to allow things to write to it, this true/false value determines if you want to use that or not

#### logkey ""
If you are using secure logging, this is the key used to allow the logger to accept logging messages from this service

#### sqsname 
This is the name of the SQS Queue. You don't need to supply the ARN (the code does a lookup for it).  This needs to match the queue name on the _dispatcher_. 

#### sqsregion 
The AWS region where your SQS Queue resides.

#### strkey 
This is the encryption key that is used to decrypt incoming Action Messages.  This needs to match the strkey in the _dispatcher_.

#### slptimeout 
The slptimeout is the sleep timeout.  This is how longer after checking the SQS Queue that it will wait before checking again.

#### worker_hosts "http://worker1.domain.com:8094|http://worker2.domain.com:8094|http://worker3.domain.com:8094"
This is a list of worker hosts (ie: hosts that will run the containers to handle actual work)  The http:// and Port numbers (:8094) are required

