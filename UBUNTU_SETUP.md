# Ubuntu 22.04 LTS Setup

## Part 1: Initialization

On your *server*, login as your administrator and create the `~/go/bin` folder which you can upload your apps to.

```bash
mkdir -p ~/go/bin
```

On your *local machine* run the following:
```bash
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/bleve-server -v
```

Finally on your *local machine* run the following (replace `xxx.xxx.xxx.xxx` with your servers IP):
```bash
rsync -avHe ssh bin/bleve-server techops@xxx.xxx.xxx.xxx:~/go/bin
```

In your *server*, verify the app works:

```bash
cd ~/go/bin
./bleve-server serve
```

If you get output then congratulations you have successfully go the app to work.

## Part 2: Volumes
If you have an attached `volume` to your `droplet` then you can save all the `bleve-server` files there!

Start by creating a default folder and setting correct permissions:

```bash
sudo mkdir /mnt/volume_tor1_01/bleve
sudo chown techops /mnt/volume_tor1_01/bleve
sudo chgrp techops /mnt/volume_tor1_01/bleve
sudo chmod u+rwx /mnt/volume_tor1_01/bleve
```

Next setup the environment variable:

```bash
export BLEVE_SERVER_HOME_DIRECTORY_PATH="/mnt/volume_tor1_01/bleve"
```

Finally run the server, create a test index and verify it exists:

```bash
cd ~/go/bin
./bleve-server serve
./bleve-server register --filename="testtenant"
ls -alh /mnt/volume_tor1_01/bleve # You should see `testtenant` when you run this command!
```

## Part 3: Linux Startup / Shutdown

### (I) Startup Script

1. Create our startup script. There are a number of ways to get environment variables for our project to work with ``systemd``. We will create a ``bash`` script which will load up the environment variables followed by loading up the application. Special thanks to https://unix.stackexchange.com/a/455283.

    ```bash
    cat > ~/go/bin/bleve-server.serve.sh
    ```

2. Copy and paste into the file.


    ```bash
    #!/bin/bash
    export BLEVE_SERVER_ADDRESS=127.0.0.1:8001
    export BLEVE_SERVER_HOME_DIRECTORY_PATH="/mnt/volume_tor1_01/bleve"
    exec /home/techops/go/bin/bleve-server serve
    ```

3. Make the script accessible.

    ```bash
    $ chmod +755 ~/go/bin/bleve-server.serve.sh
    ```

4. Run to confirm.

    ```bash
    ~/go/bin/bleve-server.serve.sh
    ```

### (II) Systemctl Integration
This section explains how to integrate our project with ``systemd`` so our operating system will handle stopping, restarting or starting.

1. While you are logged in as a ``techops`` user, please write the following into the console.

    ```bash
    sudo vi /etc/systemd/system/bleve-server.service
    ```


2. Copy and paste the following.

    ```text
    Description=Bleve Server RPC App Service
    Wants=network.target
    Requires=network.target
    After=network.target

    [Service]
    Type=simple
    DynamicUser=yes
    WorkingDirectory=/home/techops/go/bin
    ExecStart=/bin/bash /home/techops/go/bin/bleve-server.serve.sh
    Restart=always
    RestartSec=3
    SyslogIdentifier=ep-backend-serve
    StandardOutput=journal+console
    StandardError=journal+console
    User=techops
    Group=techops
    PermissionsStartOnly=True
    RuntimeDirectoryMode=0775
    ReadWritePaths=/mnt/volume_tor1_01/bleve  # This is how you grant read/write access to all files/directories inside a particular location to your application.

    [Automount]
    Where=/mnt/volume_tor1_01/bleve  # Automount forces systemctl to mount this location for our service's usage.

    [Install]
    WantedBy=multi-user.target
    ```

3. Grant access.

   ```bash
   sudo chmod 755 /etc/systemd/system/bleve-server.service
   ```


4. (Optional) If you've updated the above, you will need to run the following before proceeding.

    ```bash
    sudo systemctl daemon-reload
    ```


5. We can now start the service we created and enable it so that it starts at boot:

    ```bash
    sudo systemctl start bleve-server
    sudo systemctl enable bleve-server
    ```

6. Confirm our service is running.

    ```bash
    sudo systemctl status bleve-server
    ```

7. (OPTIONAL) Continuously see the log file in your terminal. While you keep the console open the stream of information will come in.

    ```bash
    sudo journalctl -f -u bleve-server.service
    ```

### (III) Microservice Configuration
If you have another `droplet` in the same region as this `droplet`, then we can use `private networking` and turn `bleve-server` into a microservice which has the RPC exposed.

1. Open up our startup script.

    ```bash
    cat > ~/go/bin/bleve-server.serve.sh
    ```

2. Change the IP address (ex: `xxx.xxx.xxx.xxx`) to the private address (ex: `yyy.yyy.yyy.yyy`) given to you by your virtual host provider (ex: DigitalOcean) and the port you would like to open:


    ```bash
    #!/bin/bash
    export BLEVE_SERVER_ADDRESS=yyy.yyy.yyyy.yyy:33001
    export BLEVE_SERVER_HOME_DIRECTORY_PATH="/mnt/volume_tor1_01/bleve"
    exec /home/techops/go/bin/bleve-server serve
    ```

3. Reload `systemctl` and confirm it works

    ```bash
    sudo systemctl daemon-reload
    sudo systemctl restart bleve-server
    sudo journalctl -f -u bleve-server.service
    ```

4. Your other droplet should be able to connect to this microservice using `RPC`.
