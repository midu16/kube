# kube
Mihai I. - 2020
The scope of this git is to provision an K8s distributed cluster using ansible and vagrant. 
The k8s-master contain also the ```helm-cli``` for deploying any k8s pods using helm package manager for kubernetes.

The ```jenkins-node``` role is to enable the posibility of maintaining in an CI-CD way all the custom made helm-charms. In order to access the jenkinsUI visit the following web-url ```http://192.168.50.5:8080/login``` and check the installation log for the administrative password for the first login.

## Requirements
On premise host will require the following pre-requisites :

    - Vagrant provider : 
                - VirtualBox
                - VMware
                - Hyper-V
                - Vagrant Cloud

    - Install Vagrant : https://www.vagrantup.com/downloads , recomended is to use Vagrant 2.2.9.
    
    - Install Ansible : https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html , recomended is to use ansible-playbook 2.9.7

    - InfluxDB shell version: 1.7.10

    - 
## Operating System 
OS compatibility :

    - Ubuntu 16.04 - done
    
    - Ubuntu 18.04 - ongoing

    - Debian 10 - ongoing

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

## Results
After runnning the ```$ vagrant up``` command from the working director the expected output is the following:
```
$ vagrant global-status
id       name       provider   state   directory
---------------------------------------------------------------------------
5011ca5  k8s-master virtualbox running /Users/midu/Documents/GitHub/kube
af303db  node-1     virtualbox running /Users/midu/Documents/GitHub/kube
30e3cfc  node-2     virtualbox running /Users/midu/Documents/GitHub/kube
f35b142  node-3     virtualbox running /Users/midu/Documents/GitHub/kube
8f6499f  node-4     virtualbox running /Users/midu/Documents/GitHub/kube
7c45f9d  node-5     virtualbox running /Users/midu/Documents/GitHub/kube

```
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
## Documentation
```
[1] https://helm.sh/docs/

[2] https://docs.influxdata.com/telegraf/v1.14/introduction/installation/

[3] https://www.vagrantup.com/docs

[4] https://docs.ansible.com/ansible/latest/user_guide/playbooks.html

[5] https://grafana.com/docs/

[6] https://prometheus.io

[7] https://docs.influxdata.com/telegraf/v1.14/

```
