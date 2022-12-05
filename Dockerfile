# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11
COPY build/pleiades-linux-amd64 /pleiades
ENTRYPOINT ["/pleiades"]
