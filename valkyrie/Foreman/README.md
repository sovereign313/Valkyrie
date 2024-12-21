# Valkyrie Foreman 

The Foreman tool of Valkyrie is essentially a worker orchestration tool that recieves requests, and contacts the launcher.  The main purpose of foreman is to do the
same basic job as Dispatcher, but without using AWS.  This allows organizations that don't want to use AWS SQS to keep it all in-house.  Foreman also provides 
duplication protection.

### Action Message
An Action Message is basically the bundled of information required for Valkyrie to do it's tasks. It is the information needed by the worker to actually execute whatever
task is specified.  These are the fields of that message:  
`Host`        - [Required] - The Host To Take Action On  
`Action`      - [Required] - This is basically a tool/script in a git repo to run on the specified host  
`ActionLevel` - [Required] - Currently "external" is the only acceptable option  
`Image`       - [Required] - This is the docker image to run to complete the work.  Currently "worker" is the only Valkyrie option, but you can run anything you want  
`Contact`     - [Optional] - The Person / Group to contact in the event of something going wrong (or on success if you want)  
`MiscParams`  - [Optional] - These are parameters sent to the tool/script specified in the "Action" field  
`OnRemote`    - [Optional] - This is a true/false value that specifies if the script should first be scp'd to the host and run on the host or not  

### /trigger
The /trigger end point is used to initiate a new valkyrie worker.  This can be done via CURL using just GET requests, or from some webhook.  Its entire purpose is to
actually trigger a new event/action. It takes an Action message split up into git parameters. For example:
```
curl 'http://somehost.com:8091/trigger?host=someotherhost.com&action=find_big_file&miscparams=/var&actionlevel=external&onremote=true&image=worker'
```

### /jsonbodytrigger
/jsonbodytrigger is just like /trigger, except that it reads a json response of the Action Message from a POST Body.  This is how Splunk and ELK deliver data on their
webhook triggers.  It requires the same data as in an Action Message, and can be created in Splunk with something like:

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

The above is a Splunk query that creates all the required fields in a Splunk Alert to trigger an event in Valkyrie.  The data that is being parsed is nagios perfdata (XML files)
being ingested into Splunk.

### Environment Variables
#### dupprotection true
This true/false value determines if you want to use duplication protection or not.  Basically, if it's the same Action Message (if it's the same host + action) then drop the message

#### dupprottime 30
This is how long in seconds to deny duplicated messages (in this case 30 seconds)

#### loggerurl http://localhost:8092/
This is the URL/Port of the Valkyrie Logger Service.  All of the services in Valkyrie try to contact this URL to log it's output/actions

#### usesecurelogging false
Secure logging allows the logger to require a key value to allow things to write to it, this true/false value determines if you want to use that or not

#### logkey ""
If you are using secure logging, this is the key used to allow the logger to accept logging messages from this service

#### worker_hosts "http://worker1.domain.com:8094|http://worker2.domain.com:8094|http://worker3.domain.com:8094"
This is a list of worker hosts (ie: hosts that will run the containers to handle actual work)  The http:// and Port numbers (:8094) are required
