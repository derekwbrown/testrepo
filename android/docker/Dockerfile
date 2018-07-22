FROM ubuntu:16.04
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# Environment
ENV GIMME_GO_VERSION 1.10.2
ENV GOPATH /
ENV CONDA_PATH /root/miniconda2

ENV ANDROID_NDK_VER=r17b
ENV ANDROID_SDK_VER=4333796

# RUN echo "deb http://deb.debian.org/debian wheezy contrib non-free" >> /etc/apt/sources.list
RUN apt-get -qq update && apt-get -qq install -y fakeroot curl git procps bzip2 \
  python2.7-dev python-pip \
  libsnmp-base libsnmp-dev snmp-mibs-downloader \
  build-essential pkg-config default-jdk

# RUN apt-get install -y libsystemd-journal-dev/wheezy-backports

RUN pip install --upgrade pip -i https://pypi.python.org/simple
RUN pip install -i https://pypi.python.org/simple/ setuptools

# Invoke
RUN pip install -i https://pypi.python.org/simple/ invoke
RUN pip install -i https://pypi.python.org/simple/ awscli==1.11.152
# Reno (we have to install pbr before reno)
RUN pip install -i https://pypi.python.org/simple/ pbr
RUN pip install -i https://pypi.python.org/simple/ reno

# RVM
# RUN gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
# RUN curl -sSL https://get.rvm.io | bash -s stable
# RUN /bin/bash -l -c "rvm requirements"
# RUN /bin/bash -l -c "rvm install 2.3 && rvm cleanup all"
# RUN /bin/bash -l -c "gem install bundler --no-ri --no-rdoc"

# Gimme
RUN curl -sL -o /bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
RUN chmod +x /bin/gimme
RUN gimme $GIMME_GO_VERSION
RUN apt-get install unzip
RUN eval "$(gimme)"
# Docker
RUN curl -fsSL https://get.docker.com/ | sh

# create the agent build folder within $GOPATH
RUN mkdir -p /src/github.com/DataDog/datadog-agent

# get the android ndk
RUN wget -q -O android-ndk-$ANDROID_NDK_VER-linux-x86_64.zip https://dl.google.com/android/repository/android-ndk-$ANDROID_NDK_VER-linux-x86_64.zip
# get the android SDK
RUN wget -q -O sdk-tools-linux-$ANDROID_SDK_VER.zip https://dl.google.com/android/repository/sdk-tools-linux-$ANDROID_SDK_VER.zip
RUN unzip -qq android-ndk-$ANDROID_NDK_VER-linux-x86_64.zip
RUN mkdir /android-sdk
WORKDIR /android-sdk
RUN unzip -qq /sdk-tools-linux-$ANDROID_SDK_VER.zip
WORKDIR /android-sdk/tools/bin
RUN echo y | ./sdkmanager "platforms;android-22" "platforms;android-27" > /dev/null
ENV ANDROID_HOME /android-sdk

##
## temporary; build GO 1.11 from source.  Until it's released, need to do
## this to use GO 1.11 as the compiler
##
RUN apt-get install -y vim dos2unix
RUN mkdir /go1.11
WORKDIR /go1.11
RUN git clone -b go1.11beta2 https://go.googlesource.com/go
COPY build_go111.sh /go1.11
RUN chmod +x ./build_go111.sh
RUN dos2unix build_go111.sh
RUN ./build_go111.sh

## for now, manually set GO environment
ENV GOROOT /go1.11/go
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/go1.11/go/bin
RUN go get -u github.com/golang/dep/cmd/dep
## initialize gomobile
RUN go get golang.org/x/mobile/cmd/gomobile
RUN gomobile init -ndk /android-ndk-r17b


RUN mkdir -p /src/github.com/derekwbrown
WORKDIR /src/github.com/derekwbrown
RUN git clone https://github.com/derekwbrown/testrepo
WORKDIR /src/github.com/derekwbrown/testrepo/android
RUN dep ensure
RUN chmod +x ./gradlew
RUN dos2unix gradlew
RUN ./gradlew build
RUN mkdir -p /mnttmp
RUN mkdir -p /tmp
RUN cp app/build/outputs/apk/release/app-release-unsigned.apk /tmp
