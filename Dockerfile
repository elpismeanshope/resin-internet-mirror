# AUTOGENERATED BY MAKE - DO NOT MODIFY MANUALLY
FROM resin/amd64-alpine

# Install ZIM (wikipedia etc) packages
# wiktionary_en_simple_all.zim 38M
# wikibooks_fa_all_2016-12.zim 44M
# wikibooks_ar_all_2016-12.zim 22M
#
# wikipedia_fa_all_2016-11.zim 4.5G
# wikipedia_ar_all_2016-12.zim 4.9G
# ==> 9.5G

RUN mkdir -p /usr/src/kiwix \
 && apk add --update openssl \
 && KIWIX= \
 && BIN= \
 && if [ $(uname -m) = "x86_64" ]; then \
      KIWIX=kiwix-linux-x86_64 \
      ; BIN=kiwix/bin \
    ; else \
      KIWIX=kiwix-server-arm \
      ; BIN= \
      ; \
    fi \
 && wget -O - https://download.kiwix.org/bin/$KIWIX.tar.bz2 \
      | tar -C /usr/src/kiwix -xjf - \
 && mv /usr/src/kiwix/$BIN/* /usr/bin/ \
 && rm -rf /usr/src/kiwix

# Install httrack
RUN apk add --update -t build-deps build-base zlib-dev openssl-dev \
 && wget -O - https://github.com/xroche/httrack/archive/3.48.21.tar.gz \
      | tar -C /usr/src -xzf - \
 && cd /usr/src/httrack-3.48.21 \
 && ./configure \
 && make -j8 \
 && make install \
 && apk add --update hostapd dnsmasq s6 nginx openssh py-pip python \
      usb-modeswitch \
 && sed -i 's/#PermitRootLogin.*/PermitRootLogin\ yes/' \
      /etc/ssh/sshd_config \
 && pip install ka-lite \
 && adduser -h /data/kalite kalite -D \
 && apk del build-deps \
 && rm -rf /root/.cache/ /usr/src/*

# Mirror Websites
ENV SITES http://elpissite.weebly.com
RUN mkdir -p /content/www \
 && cd /content/www \
 && httrack --ext-depth=1 --disable-security-limits --max-rate 0 $SITES \
 && find . -type f -maxdepth 1 -delete \
 && rm -rf hts-cache \
 && cd /content/www/elpissite.weebly.com \
 && for f in *.html; do sed -i 's/class=\"footer-wrap\".*/class=\"footer-wrap\" style=\"display:none;\" \/>/' $f; done \
 && for f in *.html; do sed -i 's/class=\"wsite-section-elements\".*/class=\"footer-wrap\" style=\"visibility:visible;\" \/>/' $f; done \
 && for f in *.html; do sed -i 's/scpt.parentNode.insertBefore(elem, scpt);//' $f; done \
 && for f in *.html; do sed -i 's/s.parentNode.insertBefore(ga, s);//' $f; done \
 && cd files \
 && for f in *.css; do sed -i 's/@import url(\x27http:\/\/fast.fonts.net\/t\/1.css?apiType=css&amp;projectid=b9a63dc3-765c-484e-bafe-ef372307f1b7?1485949767\x27);//' $f; done \
 && for f in *.css; do sed -i 's/http:\/\/elpissite.weebly.com\/files\/theme\/fonts/..\/elpissite.weebly.com\/files\/theme\/fonts/' $f; done \
 && cd ../../fonts.googleapis.com \
 && for f in *.css; do sed -i 's/http:\/\/fonts.gstatic.com/..\/fonts.gstatic.com/' $f; done 

# Go-Questionaire
RUN apk add -t build-deps --update go git libc-dev \
 && export GOPATH=/opt/go \
 && go get github.com/elpismeanshope/go-questionnaire \
 && go install github.com/elpismeanshope/go-questionnaire \
 && mv $GOPATH/bin/go-questionnaire /usr/bin/ \
 && rm -rf /opt/go \
 && apk del build-deps

VOLUME [ "/data" ]

COPY files /

EXPOSE 80 8080 8008

ENTRYPOINT [ "/run.sh" ]
