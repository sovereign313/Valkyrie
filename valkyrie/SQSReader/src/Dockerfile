# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.
# FROM scratch
FROM alpine
MAINTAINER Ernest E. teem III <sovereign313@gmail.com>
ADD sqsreader /sqsreader 

# Logging Stuff 
ENV loggerurl http://localhost:8092/
ENV usesecurelogging false
ENV logkey ""



# If Going Through A CNTLM Proxy
#ENV HTTP_PROXY="http://localhost:3128"
#ENV HTTPS_PROXY="http://localhost:3128"
#ENV http_proxy="http://localhost:3128"
#ENV https_proxy="http://localhost:3128"



#### AWS SPECIFIC SETTINGS ####
ENV sqsname valkyrie 
ENV sqsregion us-east-1
ENV sleeptimeout 10


# These MUST Be Changed (if using AWS)
ENV strkey your_encryption_key 
ENV AWS_ACCESS_KEY_ID <your_aws_access_key_id> 
ENV AWS_SECRET_ACCESS_KEY <your_aws_secret_access_key> 





#######  CONTAINER UPDATE STUFF  #######
RUN apk update && apk add tzdata 
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" > /etc/timezone
RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates


CMD ["/sqsreader"]


