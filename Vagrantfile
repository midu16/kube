# Vagrant OS 
IMAGE_NAME = "bento/ubuntu-16.04"
# Number of K8s worker nodes
N = 1

Vagrant.configure("2") do |config|
    config.ssh.insert_key = false

    config.vm.provider "virtualbox" do |v|
        v.memory = 1024
        v.cpus = 2
    end
      
    config.vm.define "k8s-master" do |master|
        master.vm.box = IMAGE_NAME
        master.vm.network "private_network", ip: "192.168.50.10"
        master.vm.hostname = "k8s-master"
        master.vm.provision  "playbook1", type:'ansible' do |ansible|
            ansible.playbook = "kubernetes-setup/master-playbook.yml"
	        ansible.verbose = "-vvv"
            ansible.extra_vars = {
                node_ip: "192.168.50.10",
            }
        end
    end
    (1..N).each do |i|
        config.vm.define "node-#{i}" do |node|
            node.vm.box = IMAGE_NAME
            node.vm.network "private_network", ip: "192.168.50.#{i + 10}"
            node.vm.hostname = "node-#{i}"
            node.vm.provision  "playbook1", type:'ansible' do |ansible|
                ansible.playbook = "kubernetes-setup/node-playbook.yml"
		        ansible.verbose = "-vvv"
                ansible.extra_vars = {
                    node_ip: "192.168.50.#{i + 10}",
                }
            end
        end
    end

#Configure the k8s-master helm-cli
    config.vm.define "k8s-master" do |master|
        master.vm.provision  "playbook2", type:'ansible' do |ansible|
            ansible.playbook = "helm-repo/helm-k8s-dashboard/k8s-dashboard-playbook.yml"
            ansible.verbose = "-vvv"                
            ansible.extra_vars = {
                node_ip: "192.168.50.10",
            }
        end
    end

    config.vm.define "influxDB-node" do |influx|
        influx.vm.box = IMAGE_NAME
        influx.vm.network "private_network", ip: "192.168.50.8"
        influx.vm.hostname = "influxDB-node"
        influx.vm.provision "ansible" do |ansible|
            ansible.playbook = "influx-setup/influxDB-playbook.yml"
            ansible.verbose = "-vvv"
            ansible.extra_vars = {
                node_ip: "192.168.50.8",
            }
        end
    end
    config.vm.define "jenkins-node" do |jenkins|
        jenkins.vm.box = IMAGE_NAME
        jenkins.vm.network "private_network", ip: "192.168.50.5"
        jenkins.vm.hostname = "jenkins-node"
        jenkins.vm.provision "ansible" do |ansible|
            ansible.playbook = "jenkins/jenkins-playbook.yml"
            ansible.verbose = "-vvv"
            ansible.extra_vars = {
                node_ip: "192.168.50.5",
            }
        end
    end

end
