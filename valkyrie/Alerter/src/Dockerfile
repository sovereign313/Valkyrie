# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.
# FROM scratch
FROM alpine
MAINTAINER Ernest E. Teem III <sovereign313@gmail.com>
ADD alerter /alerter

ENV loggerurl http://localhost:8092/
ENV usesecurelogging false
ENV logkey ""
ENV logfilelocation /tmp/valkyrie.log
ENV fromemail valkalerter@whatever.com

# Must Include The Port
ENV mailserver smtp.whatever.com:25

# If Your Using Twilio
ENV twilio_account AC1d10xxxxxxxxxxxxxxxxxxxxxxxxxxxx
ENV twilio_token 5e84xxxxxxxxxxxxxxxxxxxxxxxxxxxx
ENV twilio_from 1844xxxxxxx

# If Your Using AWS SNS
ENV awsregion us-east-1
ENV aws_access_key_id AKxxxxxxxxxxxxxxxxxxxxx
ENV aws_secret_key_id yF9xxxxxxxxxxxxxxxxxxxxxxxxxxxx

# If Going Through A CNTLM Proxy
#ENV HTTP_PROXY="http://localhost:3128"
#ENV HTTPS_PROXY="http://localhost:3128"
#ENV http_proxy="http://localhost:3128"
#ENV https_proxy="http://localhost:3128"


RUN apk update && apk add tzdata
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" > /etc/timezone
RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates

CMD ["/alerter"]

