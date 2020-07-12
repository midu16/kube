# kube
Mihai I. - 2020
The scope of this git is to provision an K8s distributed cluster using ansible and vagrant. 
The k8s-master contain also the ```helm-cli``` for deploying any k8s pods using helm package manager for kubernetes.


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

## Operating System 
OS compatibility :

    - Ubuntu 16.04 - done
    
    - Ubuntu 18.04 - 

    - Debian 10

## Block components
The following components are included in each VM:

    -  ```influxDB-node``` includes the following components: 

                    - influxDB

                    - telegraf - intention is that all components should have at least the ```basic metrics (cpu, disk, diskio, kernel, mem, processes, swap, system)``` exported towards influxDB and Grafana.

    - ```k8s-master``` includes the following components:
                    
                    - docker
                    
                    - telegraf

                    - kubelet 

                    - kubeadm 

                    - kubectl

    - ```node-N``` includes the following components:

                    - docker
                    
                    - telegraf

                    - kubelet 

                    - kubeadm 

                    - kubectl

On the K8s-cluster it is intended to ```helm install prometheus``` and ```helm install grafana```. Connect GrafanaUI to  prometheus and influxDB.



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
