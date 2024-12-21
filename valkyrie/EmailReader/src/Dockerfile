# This Can Be Set FROM scratch, but if you plan to use git, it's better to use alpine.
# FROM scratch
FROM alpine
MAINTAINER Ernest E. teem III <sovereign313@gmail.com>
ADD emailreader /emailreader 

ENV dupprotection true
ENV dupprottime 30

ENV loggerurl http://localhost:8092/
ENV usesecurelogging false
ENV logkey ""
ENV worker_hosts "http://worker1.domain.com:8094|http://worker2.domain.com:8094|http://worker3.domain.com:8094"

ENV mailhost mail.yourserver.tld:110
ENV mailuser youruser
ENV mailpass yourpassword
ENV mailproto pop
ENV mailusetls false

# If using imap
ENV imapfolder INBOX


# In Minutes
ENV sleeptimeout 1

# Case Sensitive
ENV trigger_subject Valkyrie_Trigger 


# If Going Through A CNTLM Proxy
#ENV HTTP_PROXY="http://localhost:3128"
#ENV HTTPS_PROXY="http://localhost:3128"
#ENV http_proxy="http://localhost:3128"
#ENV https_proxy="http://localhost:3128"


RUN apk update && apk add tzdata
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" > /etc/timezone
RUN apk del tzdata && apk add git && apk add openrc && apk add ca-certificates

CMD ["/emailreader"]

