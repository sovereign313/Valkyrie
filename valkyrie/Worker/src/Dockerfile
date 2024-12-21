# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.
# FROM scratch
FROM alpine
MAINTAINER Ernest E. teem III <sovereign313@gmail.com>
ADD worker /worker 

# Logging Stuff 
ENV loggerurl http://localhost:8092/
ENV usesecurelogging false
ENV logkey ""



# If Going Through A CNTLM Proxy
#ENV HTTP_PROXY="http://localhost:3128"
#ENV HTTPS_PROXY="http://localhost:3128"
#ENV http_proxy="http://localhost:3128"
#ENV https_proxy="http://localhost:3128"



# The Path On The Remote Host To Drop
# The Code To Run
ENV externalpath /tmp/valhalla 




# If You Want To Use A Git Repo For The Tools
ENV usegitrepo true
ENV gitrepourl https://github.com/sovereign313/valhalla.git 

# If onremote = true scp/ssh in as this user
ENV sshuser root



#######  CONTAINER UPDATE STUFF  #######
RUN apk update && apk add tzdata 
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" > /etc/timezone
RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates


CMD ["/worker"]


