package unittests

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/kubernetes-sigs/multi-tenancy/tenant/pkg/apis"
	tenancyv1alpha1 "github.com/kubernetes-sigs/multi-tenancy/tenant/pkg/apis/tenancy/v1alpha1"
	tenant2 "github.com/kubernetes-sigs/multi-tenancy/tenant/pkg/controller/tenant"
	tenantnamespace "github.com/kubernetes-sigs/multi-tenancy/tenant/pkg/controller/tenantnamespace"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var cfg *rest.Config
var c client.Client
var err error
var tenantAdminSecret = corev1.Secret{}

const timeout = time.Second * 200

var sa = &corev1.ServiceAccount{
	TypeMeta: metav1.TypeMeta{
		Kind: "ServiceAccount",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "t1-admin1",
		Namespace: "default",
	},
}

var tenant = &tenancyv1alpha1.Tenant{
	ObjectMeta: metav1.ObjectMeta{
		Name: "tenant-sample",
	},
	Spec: tenancyv1alpha1.TenantSpec{
		TenantAdminNamespaceName: "tenantadmin",
		TenantAdmins: []v1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sa.ObjectMeta.Name,
				Namespace: sa.ObjectMeta.Namespace,
			},
		},
	},
}

var tenantnamespaceObj = &tenancyv1alpha1.TenantNamespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "tenantnamespace-sample",
		Namespace: tenant.Spec.TenantAdminNamespaceName,
	},
	Spec: tenancyv1alpha1.TenantNamespaceSpec{
		Name: "t1-ns1",
	},
}

func CreateCrds() {
	tr := true
	apis.AddToScheme(scheme.Scheme)

	e := &envtest.Environment{
		CRDDirectoryPaths:  []string{filepath.Join("..", "crds")},
		UseExistingCluster: &tr,
	}

	if cfg, err = e.Start(); err != nil {
		fmt.Println(err)
	}

	e.Stop()
}

func CreateTenant(t *testing.T, g *gomega.WithT) {
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	c = mgr.GetClient()

	recFn, _ := tenant2.SetupTestReconcile(tenant2.NewReconciler(mgr))
	g.Expect(tenant2.AddManagerReconciler(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := tenant2.StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the service account
	err = c.Create(context.TODO(), sa)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Create the Tenant object and expect the tenantAdminNamespace to be created
	err = c.Create(context.TODO(), tenant)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Wait for the tenantadminnamespace to be created
	nskey := types.NamespacedName{Name: tenant.Spec.TenantAdminNamespaceName}
	adminNs := &corev1.Namespace{}
	g.Eventually(func() error { return c.Get(context.TODO(), nskey, adminNs) }, timeout).Should(gomega.Succeed())


	rolekey := types.NamespacedName{Name: "tenant-admin-role", Namespace: tenant.Spec.TenantAdminNamespaceName}
	tenantRole := &v1.Role{}
	g.Eventually(func() error { return c.Get(context.TODO(), rolekey, tenantRole) }, timeout).Should(gomega.Succeed())
	

	// userCl.Delete(context.TODO(), tenantnamespaceObj)

	// g.Eventually(requestsTenantNS, timeout).Should(gomega.Receive(gomega.Equal(expectedTenantNSRequest)))
	// We should wait until the tenantnamespace cr is deleted
	// tns := types.NamespacedName{Name: tenantnamespaceObj.Name, Namespace: tenantnamespaceObj.Namespace}
	// g.Eventually(func() error { return c.Get(context.TODO(), tns, tenantnamespaceObj) }, timeout).Should(gomega.HaveOccurred())

	// c.Delete(context.TODO(), tenant)

	saSecretName, err := tenantnamespace.FindSecretNameOfSA(c, sa.ObjectMeta.Name)
	if err != nil {
		t.Logf("Failed to get secret name of a service account: error: %+v", err)
		return
	}

	//Get secret
	saSecretKey := types.NamespacedName{Name: saSecretName, Namespace: sa.ObjectMeta.Namespace}
	tenantAdminSecret = corev1.Secret{}
	err = c.Get(context.TODO(), saSecretKey, &tenantAdminSecret)
	if err != nil {
		t.Logf("Failed to get tenant admin secret, error %+v", err)
		return
	}
}

func CreateTenantNS(t *testing.T, g *gomega.GomegaWithT) {
	//Generate user config string
	userCfgStr, err := tenantnamespace.GenerateCfgStr("kind-kind", cfg.Host, tenantAdminSecret.Data["ca.crt"], tenantAdminSecret.Data["token"], sa.ObjectMeta.Name)
	userCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(userCfgStr))
	if err != nil {
		t.Logf("failed to create user config, got an invalid object error: %v", err)
		return
	}

	//User manager and client
	mgr, err := manager.New(userCfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	userCl := mgr.GetClient()

	stopmgr, mgrStopped := tenantnamespace.StartTestManager(mgr, g)

	defer func() {
		close(stopmgr)
		mgrStopped.Wait()
	}()

	//create tenantnamespace object using user client
	
	err = userCl.Create(context.TODO(), tenantnamespaceObj)
	if err != nil {
		t.Logf("failed to create tenantnamespace object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())

	//check if tenantnamespace is created or not
	nskey := types.NamespacedName{Name: tenantnamespaceObj.Spec.Name}
	tenantNs := &corev1.Namespace{}
	g.Eventually(func() error { return c.Get(context.TODO(), nskey, tenantNs) }, timeout).
		Should(gomega.Succeed())

}

func DestroyTenant(g *gomega.WithT) {
	// Delete Tenant
	err = c.Delete(context.TODO(), tenant)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	// Delete Service Account
	err = c.Delete(context.TODO(), sa)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
