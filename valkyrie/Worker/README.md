# Valkyrie Worker

The Worker image of Valkyrie is the heart of the work. It's the final destination for an Action Message, and does all of the work specified in an Action Message.  The Worker is
an ephemeral container that dies as soon as the work is complete. The steps involved in a worker are basically to read the Action Message (which has been split out into environment
variables), check out the specific git repository (specified in the Dockerfile), scp the file (from the git repo that matches Action) to the host (if onremote=true), and actually 
run/execute the file (via SSH) or locally (if onremote=false).  Please be aware of that fact that the worker image is based on Alpine (or scratch) and will _not have an interpreter_.
So if you want to run python or bash (or any other scripting language) that the interpreter needs to be installed on the destination host, AND you'll need to set onremote=true.  We 
plan to make alternative worker images, but part of our requirment/decision is to keep the worker small, so that container launch times are fast.  


### Environment Variables

#### loggerurl http://localhost:8092/
This is the URL/Port of the Valkyrie Logger Service.  All of the services in Valkyrie try to contact this URL to log it's output/actions

#### usesecurelogging false
Secure logging allows the logger to require a key value to allow things to write to it, this true/false value determines if you want to use that or not

#### logkey ""
If you are using secure logging, this is the key used to allow the logger to accept logging messages from this service

#### externalpath /tmp/valhalla
This is the path on the remote host to put the code to run.  Since Valkyrie cleans up after itself, we like to use /tmp/valhalla, but you can have it 
scp the file to anywhere that you have permission to as the user SSHing in.

#### usegitrepo true
This is if you want to use a gitrepo to download your tools.  You don't have to do this, but the default worker image comes with no tools by default.
You can recreate the image and include the tools that you want in `externalpath + "/live/"`, but this makes much more difficult to update your toolset
outside of Valkyrie.  It is certainly recommended that you use a git repo.

#### gitrepourl https://github.com/sovereign313/valhalla.git 
This is the gitrepo to use.  You can use any repo that has the tools you want.  Just keep in mind that the Action must match the tool name in the repo

#### sshuser root
This is the user to scp/ssh into destination host as
