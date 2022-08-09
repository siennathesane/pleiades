# Use busybox as the base image
FROM busybox:musl
# Copy over the executable file
COPY build/pleiades /home/pleiades
# Run the executable file
CMD /home/pleiades
