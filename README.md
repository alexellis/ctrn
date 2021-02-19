# Add networking to containerd with CNI

Example forked by Alex Ellis and fixed up from [renatofq/ctrofb](https://github.com/renatofq/ctrofb)

## Building / testing

```sh
# Get the CNI plugins to:

# /opt/cni/bin/

# https://github.com/containernetworking/plugins

# Copy net.d/* to /etc/cni/net.d/

# cp net.d/* /etc/cni/net.d/

go build 

sudo ./ctrn create
sudo ./ctrn start

# Follow instructions to find IP and access service which is figlet on port 8080

sudo ./ctrn remove
```

## Working example end to end

### Create the container

```sh
alex@alexx:~/go/src/github.com/alexellis/ctrn$ sudo ./ctrn create
&{0xc0003b20d0 helloweb {helloweb map[] docker.io/functions/figlet:latest {io.containerd.runc.v2 <nil>} 0xc00043e000 hello-snapshot overlayfs {267854
372 63712604158 <nil>} {267854372 63712604158 <nil>} map[]}}
```

### Create the task and get an IP

```
alex@alexx:~/go/src/github.com/alexellis/ctrn$ sudo ./ctrn start
&{0xc000382410 helloweb {helloweb map[] docker.io/functions/figlet:latest {io.containerd.runc.v2 <nil>} 0xc0003f6280 hello-snapshot overlayfs {267854
372 63712604158 <nil>} {267854372 63712604158 <nil>} map[]}}
Config of interface lo: &{[0xc0003df980 0xc0003df9b0] 00:00:00:00:00:00 /proc/14082/ns/net}
Config of interface gsl0: &{[] 82:95:7e:36:62:b9 }
Config of interface vethd63267a6: &{[] 82:95:7e:36:62:b9 }
Config of interface eth0: &{[0xc0003df9e0] 9e:a5:79:af:39:ab /proc/14082/ns/net}
2019/12/22 09:36:00 Version: 0.14.4     SHA: ced4ee56dc003cf4f3baa0954ab692f4be54f57b
2019/12/22 09:36:00 Read/write timeout: 5s, 5s. Port: 8080
2019/12/22 09:36:00 Writing lock-file to: /tmp/.lock
2019/12/22 09:36:00 Metrics server. Port: 8081
```

### Find the IP for the task
```
alex@alexx:~/go/src/github.com/alexellis/ctrn$ sudo ctr task ls
TASK        PID      STATUS    
helloweb    14082    RUNNING

alex@alexx:~/go/src/github.com/alexellis/ctrn$ sudo ctr task exec --exec-id 14082 helloweb ifconfig
eth0      Link encap:Ethernet  HWaddr 9E:A5:79:AF:39:AB  
          inet addr:10.10.10.4  Bcast:10.10.10.255  Mask:255.255.255.0
          inet6 addr: fe80::9ca5:79ff:feaf:39ab/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:37 errors:0 dropped:0 overruns:0 frame:0
          TX packets:9 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0 
          RX bytes:5561 (5.4 KiB)  TX bytes:690 (690.0 B)

lo        Link encap:Local Loopback  
          inet addr:127.0.0.1  Mask:255.0.0.0
          inet6 addr: ::1/128 Scope:Host
          UP LOOPBACK RUNNING  MTU:65536  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1 
          RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
```

### Access the service from within the task

```
alex@alexx:~/go/src/github.com/alexellis/ctrn$ curl -d CNI 10.10.10.4:8080
  ____ _   _ ___ 
 / ___| \ | |_ _|
| |   |  \| || | 
| |___| |\  || | 
 \____|_| \_|___|
                 
alex@alexx:~/go/src/github.com/alexellis/ctrn$ 
```

## Firecracker notes

```bash

RUNTIME=aws.firecracker \
  CONTAINERD=/run/firecracker-containerd/containerd.sock \
  SNAPSHOTTER=devmapper \
  sudo -E ./ctrn remove

```
