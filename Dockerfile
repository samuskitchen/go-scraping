FROM alpine

RUN pip install -U docker-compose
RUN pip install -U websocket

COPY dist/scraping /bin/

EXPOSE 5001

ENTRYPOINT [ "/bin/scraping" ]