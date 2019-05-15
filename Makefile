all: depthtreed

cleanvendor:
    rm -rf ~/go/pkg/source/dep/*

depthtreed:
    go install github.com/bububa/depthtree/depthtreed

install:
    cp -r ./templates /etc/depthtreed/
    cp -f ../../../../bin/depthtreed /usr/local/bin/;
    chown root:root /usr/local/bin/depthtreed;
    chmod 755 /usr/local/bin/depthtreed;
