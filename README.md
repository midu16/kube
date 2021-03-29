# kube
Mihai I. - 2020
The scope of this git is to provision an K8s distributed cluster using ansible and vagrant. 
The k8s-master contain also the ```helm-cli``` for deploying any k8s pods using helm package manager for kubernetes.

The ```jenkins-node``` role is to enable the posibility of maintaining in an CI-CD way all the custom made helm-charms. In order to access the jenkinsUI visit the following web-url ```http://192.168.50.5:8080/login``` and check the installation log for the administrative password for the first login.

## Requirements
On premise host will require the following pre-requisites :

    - Vagrant provider : 
                - VirtualBox    - done
                - VMware        - ongoing
                - Hyper-V       - ongoing
                - Vagrant Cloud - ongoing

    - Install Vagrant : https://www.vagrantup.com/downloads , recomended is to use Vagrant 2.2.9.
    
    - Install Ansible : https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html , recomended is to use ansible-playbook 2.9.7

    - InfluxDB shell version: 1.7.10

    - bento/ubuntu-16.04 version: 202007.17.0

## Operating System 
OS compatibility :
```
    - Ubuntu 16.04  - done
    - Ubuntu 18.04  - ongoing
    - Debian 10     - ongoing
    - CentOS7       - ongoing
```
## Block components
The following components are included in each VM:

```influxDB-node``` includes the following components: 

                    - influxDB
                    - telegraf : intention is that all components should have at least the basic metrics (cpu, disk, diskio, kernel, mem, processes, swap, system) exported towards influxDB and Grafana.

```jenkins-node``` includes the following compontens:

                    - telegraf
                    - jenkins 

```k8s-master``` includes the following components:
                    
                    - docker
                    - telegraf
                    - kubelet 
                    - kubeadm 
                    - kubectl

```node-N``` includes the following components:

                    - docker
                    - telegraf
                    - kubelet 
                    - kubeadm 
                    - kubectl

On the K8s-cluster it is intended to ```helm install prometheus``` and ```helm install grafana```. Connect GrafanaUI to  prometheus and influxDB.

Architecture:

![HLD Solution Architecture](https://github.com/midu16/kube/blob/dev/architecture/architecture.png)

## Networking 
Two networks are used for the interactions of the ```workstation with the k8s-cluster```, first network is using the ```CP subnet (192.168.50.0/24)``` and the second network is the ```UP subnet (10.0.2.15/24)```. The interaction between the ```workstation with the Vagrant VMs``` is used only the  ```CP subnet (192.168.50.0/24)``` for both Control Plane and User Plane.
 

## Results
After runnning the ```$ vagrant up``` command from the working director the expected output is the following:
```
$ vagrant global-status
id       name          provider   state   directory
------------------------------------------------------------------------------
2ec5303  k8s-master    virtualbox running /Users/midu/Documents/GitHub/kube
9efce6d  node-1        virtualbox running /Users/midu/Documents/GitHub/kube
30e3cfc  node-2        virtualbox running /Users/midu/Documents/GitHub/kube
f35b142  node-3        virtualbox running /Users/midu/Documents/GitHub/kube
8f6499f  node-4        virtualbox running /Users/midu/Documents/GitHub/kube
7c45f9d  node-5        virtualbox running /Users/midu/Documents/GitHub/kube
48fd748  influxDB-node virtualbox running /Users/midu/Documents/GitHub/kube
b5c5d63  jenkins-node  virtualbox running /Users/midu/Documents/GitHub/kube
eee6504  vault-node    virtualbox running /Users/midu/Documents/GitHub/kube
```
! Note: The display restults were obtain having 1 x k8s-master and 5 x worker-node's.


Once log-in into the ```k8s-master``` we can check the status of the cluster:
```
vagrant@k8s-master:~$ kubectl get nodes
NAME         STATUS   ROLES    AGE    VERSION
k8s-master   Ready    master   24m    v1.18.6
node-1       Ready    <none>   20m    v1.18.6
node-2       Ready    <none>   17m    v1.18.6
node-3       Ready    <none>   14m    v1.18.6
node-4       Ready    <none>   11m    v1.18.6
node-5       Ready    <none>   8m8s   v1.18.6
```
Checking all the pods available after the deployment on the ```k8s-master```:
```
vagrant@k8s-master:~$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-5fbfc9dfb6-87ff2   1/1     Running   0          24m
kube-system   calico-node-42q57                          1/1     Running   0          12m
kube-system   calico-node-gwzkr                          1/1     Running   0          9m
kube-system   calico-node-kmlbk                          1/1     Running   0          24m
kube-system   calico-node-ld6t8                          1/1     Running   0          18m
kube-system   calico-node-vprp4                          1/1     Running   0          21m
kube-system   calico-node-wm47g                          1/1     Running   0          15m
kube-system   coredns-66bff467f8-fk4sc                   1/1     Running   0          24m
kube-system   coredns-66bff467f8-n6wbz                   1/1     Running   0          24m
kube-system   etcd-k8s-master                            1/1     Running   0          24m
kube-system   kube-apiserver-k8s-master                  1/1     Running   0          24m
kube-system   kube-controller-manager-k8s-master         1/1     Running   0          24m
kube-system   kube-proxy-5hmzz                           1/1     Running   0          24m
kube-system   kube-proxy-5z5dj                           1/1     Running   0          12m
kube-system   kube-proxy-g8ssj                           1/1     Running   0          15m
kube-system   kube-proxy-m9vxj                           1/1     Running   0          21m
kube-system   kube-proxy-mn2jx                           1/1     Running   0          9m
kube-system   kube-proxy-x6qk4                           1/1     Running   0          18m
kube-system   kube-scheduler-k8s-master                  1/1     Running   0          24m
```
## Improving the security
In the branch https://github.com/midu16/kube/commit/2f78bb3bf50ffcaa7e096bd5182c968067133817 the security is one of the issues during the provisioning. The following steps were determined:
```
- http to https migration                       - ongoing
- centralized encrypton password management     - ongoing
- generate certifications valid for 72h         - ongoing
- defined password choose rule                  - ongoing
```
The ```vault-node``` compnents are the following:
```
- consul
- vault
```
The back-end of the vault is consul and it can be access at http://192.168.50.6:8500/ui/dc1/services/consul/instances
The front-end is vault and it can be access at http://192.168.50.6:8200/ui/vault/secrets using the Token obtain as the output of the following command ```vault operator init```

Debugging process:
```
vagrant@vault-node:~$ sudo systemctl status vault -l
● vault.service - Vault
   Loaded: loaded (/etc/systemd/system/vault.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2020-07-19 19:32:37 UTC; 2min 24s ago
     Docs: https://www.vault.io/
 Main PID: 12903 (vault)
    Tasks: 10
   Memory: 102.8M
      CPU: 136ms
   CGroup: /system.slice/vault.service
           └─12903 /usr/bin/vault server -config=/etc/vault/config.hcl

Jul 19 19:32:37 vault-node vault[12903]:               Go Version: go1.13.12
Jul 19 19:32:37 vault-node vault[12903]:               Listener 1: tcp (addr: "0.0.0.0:8200", cluster address: "0.0.0.0:8201", max_request_duration: "1
Jul 19 19:32:37 vault-node vault[12903]:                Log Level: info
Jul 19 19:32:37 vault-node vault[12903]:                    Mlock: supported: true, enabled: true
Jul 19 19:32:37 vault-node vault[12903]:            Recovery Mode: false
Jul 19 19:32:37 vault-node vault[12903]:                  Storage: consul (HA available)
Jul 19 19:32:37 vault-node vault[12903]:                  Version: Vault v1.4.3
Jul 19 19:32:37 vault-node vault[12903]: ==> Vault server started! Log data will stream in below:
Jul 19 19:32:37 vault-node vault[12903]: 2020-07-19T19:32:37.488Z [INFO]  proxy environment: http_proxy= https_proxy= no_proxy=
Jul 19 19:32:37 vault-node vault[12903]: 2020-07-19T19:32:37.489Z [WARN]  no `api_addr` value specified in config or in VAULT_API_ADDR; falling back to
vagrant@vault-node:~$ vault status
Key                Value
---                -----
Seal Type          shamir
Initialized        false
Sealed             true
Total Shares       0
Threshold          0
Unseal Progress    0/0
Unseal Nonce       n/a
Version            n/a
HA Enabled         true
```

## K8s Helm-Charms Documentation
```
[1] https://hub.helm.sh/charts/stable/grafana/3.3.6

[2] https://hub.helm.sh/charts/stable/mysql



```
## Documentation
```
[1] https://helm.sh/docs/

[2] https://docs.influxdata.com/telegraf/v1.14/introduction/installation/

[3] https://www.vagrantup.com/docs

[4] https://docs.ansible.com/ansible/latest/user_guide/playbooks.html

[5] https://grafana.com/docs/

[6] https://prometheus.io

[7] https://docs.influxdata.com/telegraf/v1.14/

[8] https://hub.helm.sh/charts/k8s-dashboard/kubernetes-dashboard

[9] https://www.vagrantup.com/docs/networking/public_network

```

## Supporting the project
You can support the project in the following ways:

- Opening an ```Issue``` ticket. Where do you try to be comprehensive in the feature that is not working for you.

- Buy me a coffee:

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/paypalme/midu161992?locale.x=en_US)