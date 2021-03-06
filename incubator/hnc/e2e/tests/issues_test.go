package test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Issues", func() {

	const (
		nsParent = "parent"
		nsChild = "child"
		nsSub1 = "sub1"
		nsSub2 = "sub2"
		nsSub1Sub1 = "sub1-sub1"
		nsSub2Sub1 = "sub2-sub1"
		nsSubSub2 = "sub-sub2"
		nsSubChild = "sub-child"
	)

	BeforeEach(func() {
		cleanupNamespaces(nsParent, nsChild, nsSub1, nsSub2, nsSub1Sub1, nsSub2Sub1, nsSubSub2, nsSubChild)
	})

	AfterEach(func() {
		cleanupNamespaces(nsParent, nsChild, nsSub1, nsSub2, nsSub1Sub1, nsSub2Sub1, nsSubSub2, nsSubChild)
	}) 

	It("Should remove obsolete conditions CannotPropagateObject and CannotUpdateObject - issue #328", func() {
		// Setting up hierarchy with rolebinding that HNC doesn't have permission to copy.
		mustRun("kubectl create ns", nsParent)
		mustRun("kubectl create ns", nsChild)
		mustRun("kubectl hns set", nsChild, "--parent", nsParent)
		// cluster-admin is the highest-powered ClusterRole and HNC is missing some of
		// its permissions, so it cannot propagate it.
		mustRun("kubectl create rolebinding cluster-admin-rb -n", nsParent, 
			"--clusterrole='cluster-admin' --serviceaccount="+nsParent+":default")
		// Tree should show CannotPropagateObject in nsParent and CannotUpdateObject in nsChild
		runShouldContainMultiple([]string{"1) CannotPropagateObject", "2) CannotUpdateObject"}, 1, "kubectl hns tree", nsParent)
		// Remove the child and verify that the condition is gone
		mustRun("kubectl hns set", nsChild, "--root")
		// There should no longer be any conditions in parent and child
		runShouldContain("No conditions", 1, "kubectl hns describe", nsParent)
		runShouldContain("No conditions", 1, "kubectl hns describe", nsChild)
	})

	It("Should set SubnamespaceAnchorMissing condition if the anchor is missing - issue #501", func() {
		// Setting up a 3-level tree with 'parent' as the root
		mustRun("kubectl create ns", nsParent)
		// create a subnamespace without anchor by creating a full namespace with SubnamespaceOf annotation 
		mustRun("kubectl create ns", nsSub1)
		mustRun("kubectl hns set", nsSub1, "--parent", nsParent)
		mustRun("kubectl annotate ns", nsSub1, "hnc.x-k8s.io/subnamespaceOf=" + nsParent)
		// If the subnamespace doesn't allow cascadingDelete and the anchor is missing in the parent namespace, it should have 'SubnamespaceAnchorMissing' condition while its descendants shoudn't have any conditions."
		// Expected: 'sub1' namespace is not deleted and should have 'SubnamespaceAnchorMissing' condition; no other conditions."
		runShouldContain("SubnamespaceAnchorMissing", 1, "kubectl get hierarchyconfigurations.hnc.x-k8s.io -n", nsSub1, "-o yaml")
	})

	It("Should unset SubnamespaceAnchorMissing condition if the anchor is re-added - issue #501", func(){
		// set up
		mustRun("kubectl create ns", nsParent)
		// create a subnamespace without anchor by creating a full namespace with SubnamespaceOf annotation 
		mustRun("kubectl create ns", nsSub1)
		mustRun("kubectl hns set", nsSub1, "--parent", nsParent)
		mustRun("kubectl annotate ns", nsSub1, "hnc.x-k8s.io/subnamespaceOf=" + nsParent)
		runShouldContain("SubnamespaceAnchorMissing", 1, "kubectl get hierarchyconfigurations.hnc.x-k8s.io -n", nsSub1, "-o yaml")
		// If the anchor is re-added, it should unset the 'SubnamespaceAnchorMissing' condition in the subnamespace.
		// Operation: recreate the 'sub1' subns in 'parent' - kubectl hns create sub1 -n parent
		// Expected: no conditions.
		mustRun("kubectl hns create", nsSub1, "-n", nsParent)
		runShouldNotContain("SubnamespaceAnchorMissing", 1, "kubectl get hierarchyconfigurations.hnc.x-k8s.io -n", nsSub1, "-o yaml")
	})

	It("Should cascading delete immediate subnamespaces if the anchor is deleted and the subnamespace allows cascadingDelete - issue #501", func() {
		// set up
		mustRun("kubectl create ns", nsParent)
		// Creating the a branch of subnamespace
		mustRun("kubectl hns create", nsSub1, "-n", nsParent)
		mustRun("kubectl hns create", nsSub1Sub1, "-n", nsSub1)
		mustRun("kubectl hns create", nsSub2Sub1, "-n", nsSub1)
		// If the subnamespace allows cascadingDelete and the anchor is deleted, it should cascading delete all immediate subnamespaces.
		// Operation: 1) allow cascadingDelete in 'ochid1' - kubectl hns set sub1 --allowCascadingDelete=true
		// 2) delete 'sub1' subns in 'parent' - kubectl delete subns sub1 -n parent
		// Expected: 'sub1', 'sub1-sub1', 'sub2-sub1' should all be gone
		mustRun("kubectl hns set", nsSub1, "--allowCascadingDelete=true")
		mustRun("kubectl delete subns", nsSub1, "-n", nsParent)
		runShouldNotContainMultiple([]string {nsSub1, nsSub1Sub1, nsSub2Sub1}, 1, "kubectl hns tree", nsParent)
	})

	It("Should cascading delete all subnamespaces if the parent is deleted and allows cascadingDelete - issue #501", func() {
		// Setting up a 3-level tree with 'parent' as the root
		mustRun("kubectl create ns", nsParent)
		// Creating the 1st branch of subnamespace
		mustRun("kubectl hns create", nsSub1, "-n", nsParent)
		mustRun("kubectl hns create", nsSub1Sub1, "-n", nsSub1)
		mustRun("kubectl hns create", nsSub2Sub1, "-n", nsSub1)
		// Creating the 2nd branch of subnamespaces
		mustRun("kubectl hns create", nsSub2, "-n", nsParent)
		mustRun("kubectl hns create", nsSubSub2, "-n", nsSub2)
		// Creating the 3rd branch of a mix of full and subnamespaces
		mustRun("kubectl create ns", nsChild)
		mustRun("kubectl hns set", nsChild, "--parent", nsParent)
		mustRun("kubectl hns create", nsSubChild, "-n", nsChild)
		// If the parent namespace allows cascadingDelete and it's deleted, all its subnamespaces should be cascading deleted.
		// Operation: 1) allow cascadingDelete in 'parent' - kubectl hns set parent --allowCascadingDelete=true
		// 2) delete 'parent' namespace - kubectl delete ns parent
		// Expected: only 'fullchild' and 'sub-fullchild' should be left and they should have CRIT_ conditions related to missing 'parent'
		mustRun("kubectl hns set", nsParent, "--allowCascadingDelete=true")
		mustRun("kubectl delete ns", nsParent)
		mustNotRun("kubectl hns tree", nsParent)
		mustNotRun("kubectl hns tree", nsSub1)
		mustNotRun("kubectl hns tree", nsSub2)
		runShouldContain("CritParentMissing: missing parent", 1, "kubectl hns tree", nsChild)
		runShouldContain("CritAncestor", 1, "kubectl hns describe", nsSubChild)
	})
})
