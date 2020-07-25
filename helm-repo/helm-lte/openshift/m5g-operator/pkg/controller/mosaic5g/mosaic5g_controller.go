package mosaic5g

import (
	"context"
	"reflect"
	"time"

	Err "errors"

	"github.com/tig4605246/m5g-operator/internal/util"
	mosaic5gv1alpha1 "github.com/tig4605246/m5g-operator/pkg/apis/mosaic5g/v1alpha1"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_mosaic5g")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Mosaic5g Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMosaic5g{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("mosaic5g-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Mosaic5g
	err = c.Watch(&source.Kind{Type: &mosaic5gv1alpha1.Mosaic5g{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Mosaic5g
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mosaic5gv1alpha1.Mosaic5g{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMosaic5g{}


// ReconcileMosaic5g reconciles a Mosaic5g object
type ReconcileMosaic5g struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Mosaic5g object and makes changes based on the state read
// and what is in the Mosaic5g.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
// How to reconcile Mosaic5g:
// 1. Create MySQL, OAI-CN and OAI-RAN in order
// 2. If the configuration changed, restart all OAI PODs
func (r *ReconcileMosaic5g) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Mosaic5g")

	// Fetch the Mosaic5g instance
	instance := &mosaic5gv1alpha1.Mosaic5g{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue

			//Delete unused ConfigMap
			conf := r.genConfigMap(instance)
			conf.Namespace = "default"
			err = r.client.Delete(context.TODO(), conf)
			if err != nil {
				reqLogger.Error(err, "Failed to delete ")
			}
			reqLogger.Info("Mosaic5g resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Mosaic5g")
		return reconcile.Result{}, err
	}

	new := r.genConfigMap(instance)
	config := &v1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: new.GetName(), Namespace: instance.Namespace}, config)
	if err != nil && errors.IsNotFound(err) {
		// Create a configmap for cn and ran
		reqLogger.Info("Creating a new ConfigMap")
		conf := r.genConfigMap(instance)
		reqLogger.Info("conf", "content", conf)
		err = r.client.Create(context.TODO(), conf)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ConfigMap")
		}
	} else if err != nil {
		reqLogger.Error(err, "Get ConfigMap failed")
		return reconcile.Result{}, err
	}

	// Define a new MySQL deployment
	mysql := &appsv1.Deployment{}
	mysqlDeployment := r.deploymentForMySQL(instance)
	// Check if MySQL deployment already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: mysqlDeployment.GetName(), Namespace: instance.Namespace}, mysql)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", mysqlDeployment.Namespace, "Deployment.Name", mysqlDeployment.Name)
		err = r.client.Create(context.TODO(), mysqlDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mysqlDeployment.Namespace, "Deployment.Name", mysqlDeployment.Name)
			return reconcile.Result{}, err
		}
		// Define a new mysql service
		mysqlService := r.genMySQLService(instance)
		err = r.client.Create(context.TODO(), mysqlService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", mysqlService.Namespace, "Service.Name", mysqlService.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "MySQL Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Creat an oaicn deployment
	cn := &appsv1.Deployment{}
	cnDeployment := r.deploymentForCN(instance)
	// Check if the oai-cn deployment already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: cnDeployment.GetName(), Namespace: instance.Namespace}, cn)
	if err != nil && errors.IsNotFound(err) {
		if mysql.Status.ReadyReplicas == 0 {
			return reconcile.Result{Requeue: true}, Err.New("No mysql POD is ready")
		}
		reqLogger.Info("MME domain name " + instance.Spec.MmeDomainName)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", cnDeployment.Namespace, "Deployment.Name", cnDeployment.Name)
		err = r.client.Create(context.TODO(), cnDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", cnDeployment.Namespace, "Deployment.Name", cnDeployment.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully. Let's wait for it to be ready
		d, _ := time.ParseDuration("30s")
		return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
	} else if err != nil {
		reqLogger.Error(err, "CN Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Create an oaicn service
	service := &v1.Service{}
	cnService := r.genCNService(instance)
	// Check if the oai-cn service already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: cnService.GetName(), Namespace: instance.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		err = r.client.Create(context.TODO(), cnService)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", cnService.Namespace, "Service.Name", cnService.Name)
			return reconcile.Result{}, err
		}
	}
	// Create an oairan deployment
	ran := &appsv1.Deployment{}
	ranDeployment := r.deploymentForRAN(instance)
	// Check if the oai-ran deployment already exists, if not create a new one
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ranDeployment.GetName(), Namespace: instance.Namespace}, ran)
	if err != nil && errors.IsNotFound(err) {
		if cn.Status.ReadyReplicas == 0 {
			d, _ := time.ParseDuration("10s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, Err.New("No oai-cn POD is ready, 10 seconds backoff")
		}
		reqLogger.Info("Sheeps are ready")
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", ranDeployment.Namespace, "Deployment.Name", ranDeployment.Name)
		err = r.client.Create(context.TODO(), ranDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", ranDeployment.Namespace, "Deployment.Name", ranDeployment.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "RAN Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := instance.Spec.Size
	if *cn.Spec.Replicas != size {
		cn.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), cn)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cn.Namespace, "Deployment.Name", cn.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the Mosaic5g status with the pod names
	// List the pods for this instance's deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(util.LabelsForMosaic5g(instance.GetName()))
	listOps := &client.ListOptions{Namespace: instance.Namespace, LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods", "Mosaic5g.Namespace", instance.Namespace, "Mosaic5g.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := util.GetPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update Mosaic5g status")
			return reconcile.Result{}, err
		}
	}

	// Check configmap is fine or not. If it's changed, update ConfigMap and restart cn ran
	if err == nil {
		if reflect.DeepEqual(new.Data, config.Data) {
			reqLogger.Info("newConf equals config")
		} else {
			reqLogger.Info("newConf does not equals config")
			reqLogger.Info("Update ConfigMap and deleting CN and RAN")
			err = r.client.Update(context.TODO(), new)
			//Should only kill the POD
			err = r.client.Delete(context.TODO(), cnDeployment)
			err = r.client.Delete(context.TODO(), ranDeployment)
			// Spec updated - return and requeue
			d, _ := time.ParseDuration("10s")
			return reconcile.Result{Requeue: true, RequeueAfter: d}, nil
		}

	}

	// Everything is fine, Reconcile ends
	return reconcile.Result{}, nil
}

// deploymentForCN returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForCN(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	cnName := m.Spec.MmeDomainName
	//ls := util.LabelsForMosaic5g(m.Name + cnName)
	replicas := m.Spec.Size
	labels := make(map[string]string)
	labels["app"] = "oaicn"
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/oaicn"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + cnName,
			Namespace:   m.Namespace,
			Labels:      labels,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           m.Spec.CNImage,
						Name:            "oaicn",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "mosaic5g-cn",
						}},
					}},
					Affinity: util.GenAffinity("cn"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}},
					},
				},
			},
		},
	}
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForRAN returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForRAN(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	ls := util.LabelsForMosaic5g(m.Name)
	replicas := m.Spec.Size
	Annotations := make(map[string]string)
	Annotations["container.apparmor.security.beta.kubernetes.io/"+m.Name+"-"+"oairan"] = "unconfined"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.GetName() + "-" + "oairan",
			Namespace:   m.Namespace,
			Annotations: Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           m.Spec.RANImage,
						Name:            "oairan",
						Command:         []string{"/sbin/init"},
						SecurityContext: &corev1.SecurityContext{Privileged: util.NewTrue()},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "cgroup",
							ReadOnly:  true,
							MountPath: "/sys/fs/cgroup/",
						}, {
							Name:      "module",
							ReadOnly:  true,
							MountPath: "/lib/modules/",
						}, {
							Name:      "usrp",
							ReadOnly:  true,
							MountPath: "/dev/bus/usb/",
						}, {
							Name:      "mosaic5g-config",
							MountPath: "/root/config",
						}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "mosaic5g-ran",
						}},
					}},
					Affinity: util.GenAffinity("ran"),
					Volumes: []corev1.Volume{{
						Name: "cgroup",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/sys/fs/cgroup/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "module",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/lib/modules/",
								Type: util.NewHostPathType("Directory"),
							},
						}}, {
						Name: "mosaic5g-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: "mosaic5g-config"},
							},
						}}, {
						Name: "usrp",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/dev/bus/usb/",
								Type: util.NewHostPathType("Directory"),
							},
						}},
					},
				},
			},
		},
	}
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// deploymentForMySQL returns a Core Network Deployment object
func (r *ReconcileMosaic5g) deploymentForMySQL(m *mosaic5gv1alpha1.Mosaic5g) *appsv1.Deployment {
	//ls := util.LabelsForMosaic5g(m.Name + cnName)
	var replicas int32
	replicas = 1
	labels := make(map[string]string)
	labels["app"] = "oai"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql",
			Namespace: m.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "mysql:5.6",
						Name:  "mysql",
						Env: []corev1.EnvVar{
							{Name: "MYSQL_ROOT_PASSWORD", Value: "linux"},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
					}},
					Affinity: util.GenAffinity("cn"),
				},
			},
		},
	}
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// genConfigMap will generate a configmap from ReconcileMosaic5g's spec
func (r *ReconcileMosaic5g) genConfigMap(m *mosaic5gv1alpha1.Mosaic5g) *v1.ConfigMap {
	genLogger := log.WithValues("Mosaic5g", "genConfigMap")
	// Make specs into map[name][value]
	datas := make(map[string]string)
	d, err := yaml.Marshal(&m.Spec)
	if err != nil {
		log.Error(err, "Marshal fail")
	}
	datas["conf.yaml"] = string(d)
	cm := v1.ConfigMap{
		Data: datas,
	}
	cm.Name = "mosaic5g-config"
	cm.Namespace = m.Namespace
	genLogger.Info("Done")
	return &cm
}

// genCNService will generate a service for oaicn
func (r *ReconcileMosaic5g) genCNService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oaicn"
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "enb", Port: 2152},
			{Name: "hss-1", Port: 3868},
			{Name: "hss-2", Port: 5868},
			{Name: "mme", Port: 2123},
			{Name: "spgw-1", Port: 3870},
			{Name: "spgw-2", Port: 5870},
		},
		Selector:  selectMap,
		ClusterIP: "None",
	}
	service.Name = "oaicn"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}

// genMySQLService will generate a service for oaicn
func (r *ReconcileMosaic5g) genMySQLService(m *mosaic5gv1alpha1.Mosaic5g) *v1.Service {
	var service *v1.Service
	selectMap := make(map[string]string)
	selectMap["app"] = "oai"
	service = &v1.Service{}
	service.Spec = v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{Name: "mysql-port", Port: 3306},
		},
		Selector:  selectMap,
		ClusterIP: "None",
	}
	service.Name = "mysql"
	service.Namespace = m.Namespace
	// Set Mosaic5g instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}
