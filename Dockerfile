FROM alpine:3.10

WORKDIR /home

ADD util_install.sh util_install.sh
RUN sh util_install.sh && rm util_install.sh

ADD app app
ADD util_dump.sh util_dump.sh
ADD util_run.sh util_run.sh

CMD ["sh", "util_run.sh"]