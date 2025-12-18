#!/bin/bash

# Create the user if it doesn't exist
if ! id -u user > /dev/null 2>&1; then
    useradd -m -s /bin/bash user
fi

# Set the Samba password for the user
echo -e "password\npassword" | smbpasswd -a -s user
smbpasswd -e user

# Start Samba services
smbd --foreground --no-process-group