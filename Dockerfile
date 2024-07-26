FROM ubuntu:22.04

# Set the working directory
RUN mkdir /app
WORKDIR /app

# Copy the built binary to the working directory
COPY ./bin/olsrv_bin /app
# Copy all files in content in too
COPY ./content /app/content

# Expose the port the application runs on
EXPOSE 8080

# Set the initial command to run
CMD ["/app/olsrv_bin", "-dbn", "/olsrv_vol/olsrv_db.db"]