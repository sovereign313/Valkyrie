# Valkyrie Email Reader 

The Email Reader of Valkyrie is just another method of allowing inputs into Valkyrie.  We considered having the foreman be a pluginnable interface, but we decided on the 
old *nix motto of `do one thing, and do it well`, as opposed to bloating the other services.  The Email Reader allows you to specify a POP or IMAP email server and account
to login to, which it will then look for a specific subject.  If the subject is found, then it looks at the body to create an Action Message.  If the body is malformed, 
Valkyrie will simply ignore the message, log that it was malformed, and move on.
 

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

#### mailhost mail.yourserver.tld:110
This is the email server to connect to in order to get email messages

#### mailuser youruser
This is the email user account to login with

#### mailpass yourpassword
The email password to use for the account

#### mailproto pop
If you want to use pop/pop3 or imap

#### mailusetls false
If you want to use TLS/SSL for your communications

#### trigger_subject Valkyrie_Trigger
This is the subject that will trigger Valkyrie.  This is case sensitive, and if the email does not match exactly, Valkyrie will ignore the message.

#### imapfolder INBOX
The folder (if using IMAP) to check for Valkyrie Messages

#### sleeptimeout 1
How long to sleep between checks to the email server.  If you are using something like Gmail, you may want to check their policy on time between checks.
Gmail's policy is roughly one check every 10 minutes.

