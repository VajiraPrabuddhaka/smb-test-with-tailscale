FROM ubuntu:22.04

# Install Samba and utilities
RUN apt-get update && \
    apt-get install -y samba samba-common-bin && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create shared directory
RUN mkdir -p /shared/public

# Copy configuration files
COPY smb.conf /etc/samba/smb.conf
COPY startup.sh /startup.sh
RUN chmod +x /startup.sh

# Expose SMB ports
EXPOSE 139 445

# Start Samba
CMD ["/startup.sh"]
