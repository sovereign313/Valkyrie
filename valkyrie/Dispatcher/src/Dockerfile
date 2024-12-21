FROM alpine
MAINTAINER Ernest E. Teem III <sovereign313@gmail.com>
ADD dispatcher dispatcher

# These CAN Be Changed (If You Want)
ENV sqsname valkyrie 
ENV sqsregion us-east-1
ENV dupprotection true
ENV dupprottime 30

# If Using A CNTLM Proxy
#ENV HTTP_PROXY="http://localhost:3128"
#ENV HTTPS_PROXY="http://localhost:3128"
#ENV http_proxy="http://localhost:3128"
#ENV https_proxy="http://localhost:3128"


# These MUST Be Changed
ENV strkey your_encryption_key 
ENV AWS_ACCESS_KEY_ID <your_aws_access_key_id>
ENV AWS_SECRET_ACCESS_KEY <your_aws_secret_access_key> 


RUN apk update
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/US/Eastern /etc/localtime
RUN echo "US/Eastern" > /etc/timezone
RUN apk del tzdata
RUN apk add git 
RUN apk add openrc
RUN apk add ca-certificates


EXPOSE 80
ENTRYPOINT ["/dispatcher"]
