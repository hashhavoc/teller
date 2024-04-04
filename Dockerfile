FROM scratch
COPY teller /usr/bin/teller
ENTRYPOINT ["/usr/bin/teller"]