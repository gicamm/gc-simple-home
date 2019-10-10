#----
# Build stage
#----
FROM golang:1.12.4-alpine3.9 as buildstage
WORKDIR /$GOPATH/src/github.com/platinasystems/tiles

# Install git
RUN apk --update add git

# Enable go modules
ENV GO111MODULE on

# Enable access to private github repositories
# Token is the "PlatinaBuilder" one related to "giovannimorana" github account
ENV GITHUB_TOKEN a556328a954548f9b3f8e5bf2bc2b9d84a30c479
RUN git config --global url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Populate the module cache based on the go.{mod,sum}
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Install build required packages
RUN apk --update add make gcc tar rsync

# Download collector go.{mod,sum} for collector
WORKDIR /$GOPATH/src/github.com/platinasystems
RUN git clone https://github.com/platinasystems/go-collector-system
WORKDIR go-collector-system
# Go on specific commit
RUN git fetch && git reset --hard aa2d2f271278789e9dfbd1717e7e48b889571e29
RUN go mod download
RUN make build distMode=dir ignoreMissing=yes
RUN COLLECTOR_DIR=$(ls dist/ | grep system) && cp dist/$COLLECTOR_DIR/system-collector /root/system-collector

# Download turnkey-kubespray
WORKDIR /$GOPATH/src/github.com/platinasystems
RUN git clone https://github.com/platinasystems/turnkey-kubespray
WORKDIR turnkey-kubespray
# Go on specific commit
RUN git fetch && git reset --hard 8febfc7fd450d09f08509a1ee1f5bc16853fa379
RUN cp -r /$GOPATH/src/github.com/platinasystems/turnkey-kubespray /root/

# Download ops repo
WORKDIR /$GOPATH/src/github.com/platinasystems
RUN git clone https://github.com/platinasystems/ops.git
# Go on specific commit
WORKDIR ops
RUN git fetch && git reset --hard c9f9f33de633405e1c46caecb2c7c5ddc8b3c10e
RUN cp -r /$GOPATH/src/github.com/platinasystems/ops /root/

#Build pccserver
WORKDIR /$GOPATH/src/github.com/platinasystems/tiles
COPY . .
RUN make build distMode=dir ignoreMissing=yes execMode=debug collectorPath=/root/system-collector turnkeyKubesprayPath=/root/turnkey-kubespray opsPath=/root/ops
RUN mv dist/* /

# Download common docker build scripts
WORKDIR /$GOPATH/src/github.com
RUN git clone https://github.com/platinasystems/installation
WORKDIR installation
RUN git reset --hard 8c9703fc6303622d3e046e794c2e871b5b52f094
RUN cp -r product/docker/common/scripts /scripts

#----
# Microservice stage
#----
FROM debian:8
#Install required packages
RUN echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" >> /etc/apt/sources.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367
RUN apt-get update && apt-get install -y postgresql-client=9.4+165+deb8u3 python-pip xz-utils less rsync bash python-netaddr git supervisor
# Install pip
RUN easy_install -U pip
#Configure installed packages
# Set rsync and ssh symlink with path searched by ansible
RUN ln -s /usr/bin/rsync /usr/sbin/rsync
RUN ln -s /usr/bin/ssh /usr/sbin/ssh

#Copy microservice data
COPY --from=buildstage  /pccserver-* /home/
RUN MS_DIR=$(ls /home | grep pccserver) && export MS_DIR=$MS_DIR
WORKDIR /home/$MS_DIR

#Install ansible
COPY build/docker/ansible/ansible_2.7.10-1.deb .
RUN dpkg -i ansible_2.7.10-1.deb || apt-get -f install -y
RUN rm ansible_2.7.10-1.deb

#Update hosts.yml
COPY build/docker/ansible/hosts.yml ansible/infra/inventory/sites

#Install db setup common scripts
RUN mkdir -p scripts
COPY --from=buildstage /scripts scripts
RUN chmod +x scripts/microservice/*

# Copy schema setup script
COPY pccserver/sql/schema.sql scripts/db/

#Setup debug env if enabled
ARG ENABLE_DEBUG
RUN chmod u+x scripts/debug/debug-env-setup.sh && scripts/debug/debug-env-setup.sh -l go -d $ENABLE_DEBUG -i debian

#Setup kubespray repo
WORKDIR turnkey-kubespray
RUN bash setup_repo.sh

# Set default variables
ENV ANSIBLE_VERBOSITY=0
#Run container
WORKDIR /home/$MS_DIR

##
## Allow tunneling
##
RUN echo "deb http://archive.canonical.com/ubuntu bionic partner" >> /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y curl net-tools bash-completion iproute2 bridge-utils nano kmod dbus

COPY pccserver/ansible/infra/network/tunnel/scripts/tunnel.sh scripts/tunnel/

RUN apt-get install -y openssh-server
RUN mkdir /var/run/sshd

RUN echo 'root:root' |chpasswd

RUN sed -ri 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -ri 's/^#?GatewayPorts\s+.*/GatewayPorts yes/' /etc/ssh/sshd_config
RUN sed -ri 's/^#?PermitTunnel\s+.*/PermitTunnel yes/' /etc/ssh/sshd_config
RUN sed -ri 's/^#?TCPKeepAlive\s+.*/TCPKeepAlive yes/' /etc/ssh/sshd_config
RUN sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config

RUN echo "PermitTunnel yes" >> /etc/ssh/sshd_config
RUN echo "AllowTCPForwarding yes" >> /etc/ssh/sshd_config

#Allow IP forwarding
RUN sed -ir 's/#{1,}?net.ipv4.ip_forward ?= ?(0|1)/net.ipv4.ip_forward = 1/g' /etc/sysctl.conf
RUN mkdir /root/.ssh
RUN systemctl enable ssh
RUN service ssh restart


#Run with supervisord
COPY build/docker/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]
