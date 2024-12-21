# Valkyrie Dispatcher

The Dispatcher tool of Valkyrie is used to allow an analyzing engine (Splunk, ELK, etc) or a custom creation (some web page, curl, etc) to trigger events that 
automatically take action on a remote host.  The Dispatcher works by taking the data via a web api (it's not a proper RESTful API), parsing it and putting
the message in AWS SQS.

The Dispatcher has 2 trigger calls, a bunch of setter calls, and some misc calls.

### Trigger Calls

##### Trigger Call (Generic)
The generic trigger call /trigger requires 3 values sent via a GET call.  It requires the Host (the host/ip to take action on), the action to take, the Contact 
(who to tell that we took action), and optionally Params, a misc set of parameters (additional parameters that get passed to the workers).  This data is encrypted prior to being 
sent to SQS. This is a GET request.


##### JSONBodyTrigger Call
The splunk trigger call /jsonbodytrigger also requires 3 values sent as a JSON body in a POST request.  It requires dsthost (the host/ip to take action on), the dstfn 
(the action to run on the remote host), the dstcontact (who to tell that we took action), and optionally dstparams, a misc set of parameters
(additional parameters that get passed to the workers).  This data is encrypted prior to being sent to SQS.  The values *must* be in the 'result' parent to the JSON.
If you're using splunk, this is the default behavior anyway.  You can create these variables to be sent using eval in Splunk. This is a POST request


### Setter Calls

##### Set SQS Region
The setter call /setsqsregion is a GET request that requires one parameter: sqsregion.  This will set the AWS region your SQS.  You can optionally set this as
an environment variable prior to running the Dispatcher daemon.  The default value set for sqsregion is: us-east-2. 

##### Set SQS Name
The setter call /setsqsname is a GET request that requires one parameter: sqsname.  This will set the name of your AWS SQS Queue.  You can optionally set this as
an environment variable prior to running the Dispatcher daemon.  The default value set for sqsname is: dispatcher.

##### Set SNS Topic
The setter call /setsnstopic (note the SNS) is a GET request that requires one parameter: snstopic.  This will set the name of your AWS SNS Topic.  You can optionally set
this as an environment variable prior to running the Dispatcher daemon.  The default value set for snstopic is: dispatcher.
This is used to send a notification when it is unable to write to the log.  If this fails, it will attempt SMTP Email.

##### Set SNS Topic ARN
The setter call /setsnstopicarn (note the SNS) is a GET request that requires one parameter: snstopicarn.  Thhis will set the ARN of your AWS SNS Topic.  You can optionally set
this as an environment variable prior to running the Dispatcher daemon.  *There is no default value set for snstopicarn*.  If this is left blank, Dispatcher will try to send
the notification (of a failed write to the log file) via SMTP.

##### Set Encrypt Key
The setter call /setstrkey is a GET request that requires one parameter: strkey.  This will set the encryption key that is used to encrypt messages being put on AWS SQS.  You can 
optionally set this as an environment variable prior to running the Dispatcher daemon.  *The choice was on purpose to have no default* and the Dispatcher will simply refuse to
send messages over the interwebz without one set.

##### Set To Email
The setter call /settoemail is a GET request that requires on parameter: toemail.  This will set the email address to send to.  You can optionally set this as an environment variable
prior to running the Dispatcher daemon.  *There is no default value set for toemail*.  If this is blank, Dispatcher will not be able to send emails on failed log writes.

##### Set From Email
The setter call /setfromemail is a GET request that requires one parameter: fromemail.  This will set the email address of the sender.  You can optionally set this as an environment variable
prior to running the Dispatcher daemon.  *There is no default value set for fromemail*.  If this is blank, Dispatcher will not be able to send emails on failed log writes.

##### Set Email Server
The setter call /setemailserver is a GET request that requires one parameter: emailserver.  This will set the email server to send email to.  You can optionally set this as an environment variable
prior to running the Dispatcher daemon.  *There is no default value set for emailserver*.  If this is blank, Dispatcher will not be able to send emails on failed log writes.

##### Set AWS Access Key ID
The setter call /setawsaccesskey is a GET request that requires one parameter: awsaccesskey.  This will set the Access Key ID for your AWS Credentials.  It literally just sets the environment variable
that the AWS API already will check for.  It's completely optional, but can easily take the place of a shared credentials file (~/.aws/credentials).  Obviously, this must be set
prior to making calls to AWS or the call will fail (unless you have set elsewhere).

##### Set AWS Secret Key
The setter call /setawssecretkey is a GET request that requires one parameter: awssecretkey.  This will set the Secret Key for your AWS Credentials.  It literally just sets the environment variable
that the AWS API already will check for.  It's completely optional, but can easily take the place of a shared credentials file (~/.aws/credentials).  Obviously, this must be set
prior to making calls to AWS or the call will fail (unless you have set elsewhere).

##### Set Duplication Protection
The setter call /setdupprotection is a GET request that requires one parameter: dupprotection.  The default value for this is TRUE/ON.  This will set the state of duplication protection of messages.  
If enabled, this will check for an epoch time stamp and discard any messages that match both the Host and Action if received within a configurable amount of time (default 30 mins).  This is 
configurable on the fly, and can be set with an environment variable also.

##### Set Duplication Protection Time
The setter call /setdupprottime is a GET request that requires one parameter: dupprottime.  This will set the time in minutes that messages matching both Host and Action will be discarded if 
received in a shorter span of time than dupprottime (default 30 mins).  This value only has an impact if dupprotection is enabled.  This is configurable on the fly, and can be set with an 
environment variable also.

### Misc Calls


##### Echo 
The echo call /echo will simply return whatever POST data was sent to /jsonbodytrigger in the body _the last time /jsonbodytrigger was called_.  
This is used to view the Splunk JSON data that was sent (so that you can verify the required parameters (dsthost, dstfn, dstcontact, and optionally dstparams).
This will *not* echo your current call.  This is a GET request.

##### Ping
The ping call /ping will return "pong".  This is simply used to verify the service is up.

##### Who Are You
The whoareyou call /whoareyou will return "dispatcher".  This is/can be used service identification.


## Examples 


##### Example of Splunk Query 
```
index="nagios" source=*/production* SERVICEOUTPUT mailq NOT *empty*
| replace /usr/local/nagios/share/perfdata/* with * in source
| eval dsthost=source
| eval dstfn=source
| eval dstfnlevel="external"
| eval dstrmpriority="delayed"
| eval dstusesns="true"
| eval dstcontact="me@myplace.com"
| replace */Mail_Queue.xml with * in dsthost
| replace */Mail_Queue.xml with Mailq in dstfn
| table dsthost dstfn dstcontact
| dedup dsthost dstfn dstcontact
```
This will have splunk parse the nagios perfdata for production servers.  It will parse the hostname, and the function to take action on.  it then sets the function level
(what type of activity the action is... either an internal function, an external script, a worker plugin).  It then sets the queue removal priority (delete from SQS before
the work is done or after), if we should SNS to notify on work completion, and the person to contact.
