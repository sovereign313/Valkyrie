# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.
# FROM scratch
FROM alpine
MAINTAINER Ernest E. teem III <sovereign313@gmail.com>
ADD logger /logger

# These CAN Be Changed (If You Want)
ENV useeventstreams false 
ENV eshostport splunk.yourdomain.com:5020



# Change this to require a key for logging
# And provide the key used
ENV usesecurelogging false
ENV logkey "" 
ENV logfilelocation /tmp/valkyrie.log

# If Going Through A CNTLM Proxy
#ENV HTTP_PROXY="http://localhost:3128"
#ENV HTTPS_PROXY="http://localhost:3128"
#ENV http_proxy="http://localhost:3128"
#ENV https_proxy="http://localhost:3128"


RUN apk update && apk add tzdata
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" > /etc/timezone
RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates

CMD ["/logger"]

