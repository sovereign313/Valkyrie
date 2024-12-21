Valkyrie is a Simple Lightweight Infrastructure Automation tool.  It's a collection of tools that allow an analyzing engine (ELK, Splunk, etc) to take action on remote systems
(or on your own, if you're into that kinda thing).  The idea is that something like Splunk can trigger a webhook (to the Dispatcher/Foreman Web API). It then triggers a worker
(either via AWS SQS [dispatcher] or direct URL call [foreman]) that checks out a git repo with a list of tools, and runs the specific tool.

For example: Splunk detects that a disk is filling up, triggers it's alert action (which hits the foreman/dispatcher).  The Foreman triggers an event in the launcher (or the 
dispatcher inserts the message into SQS, and the SQS Reader grabs the SQS Message, and triggers the Launcher).  The Launcher spins up a new worker container, which checks out
a git repo, sshes into the host, and cleans up large log files (or connects to vmware, to add more disk, and automatically adds the disk to the LVM, etc).

This is all done without admin intervention.  

Valkyrie also has reporting tools.  For example, when a CPU Usage alert is triggered, it can SSH into the server, and report back which PID is wreaking havok on your system.
Currently it uses ssh keys to access the servers, and runs as root.  This isn't required, but valkyrie is much more effective when run as root. We plan on providing access for
a common password for hosts, or for obtaining the passwords from [Hashicorp Vault](https://www.vaultproject.io/) but beefing up the toolset of [Valhalla](https://github.com/sovereign313/valhalla)
(the default repo for Valkyrie) is more of a priority right now.

This is Valkyrie.  Slayer of infrastructure burdens.

