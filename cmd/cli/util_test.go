package main

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tassert "github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openservicemesh/osm/pkg/constants"
)

func TestAnnotateErrorMessageWithActionableMessage(t *testing.T) {
	type test struct {
		errorMsg     string
		name         string
		namespace    string
		exceptionMsg string
		annotatedMsg string
	}

	actionableMsg := "Use flags to modify the command to suit your needs"

	testCases := []test{
		{
			"Error message with args such as [name: %s], [namespace: %s], and [err: %s]",
			"osm-name",
			"osm-namespace",
			"osm-exception",
			"Error message with args such as [name: osm-name], [namespace: osm-namespace], and [err: osm-exception]\n\n" + actionableMsg,
		},
	}

	for _, tc := range testCases {
		t.Run("Testing annotated error message", func(t *testing.T) {
			assert := tassert.New(t)

			assert.Equal(
				tc.annotatedMsg,
				annotateErrorMessageWithActionableMessage(actionableMsg, tc.errorMsg, tc.name, tc.namespace, tc.exceptionMsg).Error())
		})
	}
}

func TestAnnotateErrorMessageWithOsmNamespace(t *testing.T) {
	type test struct {
		errorMsg     string
		name         string
		namespace    string
		exceptionMsg string
		annotatedMsg string
	}

	osmNamespaceErrorMsg := fmt.Sprintf(
		"Note: The command failed when run in the OSM namespace [%s].\n"+
			"Use the global flag --osm-namespace if [%s] is not the intended OSM namespace.",
		settings.Namespace(), settings.Namespace())

	testCases := []test{
		{
			"Error message with args such as [name: %s], [namespace: %s], and [err: %s]",
			"osm-name",
			"osm-namespace",
			"osm-exception",
			"Error message with args such as [name: osm-name], [namespace: osm-namespace], and [err: osm-exception]\n\n" + osmNamespaceErrorMsg,
		},
	}

	for _, tc := range testCases {
		t.Run("Testing annotated error message", func(t *testing.T) {
			assert := tassert.New(t)

			assert.Equal(
				tc.annotatedMsg,
				annotateErrorMessageWithOsmNamespace(tc.errorMsg, tc.name, tc.namespace, tc.exceptionMsg).Error())
		})
	}
}

var _ = Describe("Test getting pretty printed output of a list of meshes", func() {
	var (
		meshInfoList []meshInfo
	)

	Context("empty mesh list", func() {

		meshInfoList = []meshInfo{}
		pp := getPrettyPrintedMeshInfoList(meshInfoList)

		It("should have correct output", func() {
			Expect(pp).To(Equal("\nMESH NAME\tMESH NAMESPACE\tVERSION\tADDED NAMESPACES\n"))
		})
	})

	Context("non-empty mesh list", func() {

		meshInfoList = []meshInfo{
			{
				name:                "m1",
				namespace:           "ns1",
				version:             "v1",
				monitoredNamespaces: []string{"mn1", "mn2", "mn3"},
			},
			{
				name:                "m2",
				namespace:           "ns2",
				version:             "v2",
				monitoredNamespaces: []string{"mn4", "mn5", "mn6"},
			},
		}

		It("should have correct output", func() {
			Expect(getPrettyPrintedMeshInfoList(meshInfoList)).To(Equal("\nMESH NAME\tMESH NAMESPACE\tVERSION\tADDED NAMESPACES\nm1\tns1\tv1\tmn1,mn2,mn3\nm2\tns2\tv2\tmn4,mn5,mn6\n"))
		})

	})
})

var _ = Describe("Test getting pretty printed output of smi info of a list of meshes", func() {
	var (
		meshSmiInfoList []meshSmiInfo
	)

	Context("empty mesh list", func() {
		meshSmiInfoList = []meshSmiInfo{}
		pp := getPrettyPrintedMeshSmiInfoList(meshSmiInfoList)

		It("should have correct output", func() {
			Expect(pp).To(Equal("\nMESH NAME\tMESH NAMESPACE\tSMI SUPPORTED\n"))
		})
	})

	Context("non-empty mesh list", func() {
		meshSmiInfoList = []meshSmiInfo{
			{
				name:                 "m1",
				namespace:            "ns1",
				smiSupportedVersions: []string{"smi1", "smi2", "smi3"},
			},
			{
				name:                 "m2",
				namespace:            "ns2",
				smiSupportedVersions: []string{"smi4", "smi5", "smi6"},
			},
		}

		pp := getPrettyPrintedMeshSmiInfoList(meshSmiInfoList)

		It("should have correct output", func() {
			Expect(pp).To(Equal("\nMESH NAME\tMESH NAMESPACE\tSMI SUPPORTED\nm1\tns1\tsmi1,smi2,smi3\nm2\tns2\tsmi4,smi5,smi6\n"))
		})
	})
})

// helper function for tests that adds deployment to the clientset
func addDeployment(fakeClientSet kubernetes.Interface, depName string, meshName string, namespace string, osmVersion string, isMesh bool) (*v1.Deployment, error) {
	dep := createDeployment(depName, meshName, osmVersion, isMesh)
	return fakeClientSet.AppsV1().Deployments(namespace).Create(context.TODO(), dep, metav1.CreateOptions{})
}

// helper function for tests that creates a deployment for mesh and non-mesh deployments
func createDeployment(deploymentName, meshName string, osmVersion string, isMesh bool) *v1.Deployment {
	labelMap := make(map[string]string)
	if isMesh {
		labelMap[constants.AppLabel] = constants.OSMControllerName
		labelMap["meshName"] = meshName
		labelMap[constants.OSMAppVersionLabelKey] = osmVersion
	} else {
		labelMap[constants.AppLabel] = "non-mesh-app"
	}
	dep := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deploymentName,
			Labels: labelMap,
		},
	}
	return dep
}

func TestNamespacedNameFrom(t *testing.T) {
	testCases := []struct {
		name      string
		in        string
		out       types.NamespacedName
		expectErr bool
	}{
		{
			name:      "valid namespaced name",
			in:        "foo/bar",
			out:       types.NamespacedName{Namespace: "foo", Name: "bar"},
			expectErr: false,
		},
		{
			name:      "invalid namespaced name with no separator",
			in:        "foobar",
			out:       types.NamespacedName{},
			expectErr: true,
		},
		{
			name:      "invalid namespaced name with multiple separators",
			in:        "foo/bar/baz",
			out:       types.NamespacedName{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := tassert.New(t)

			actual, err := namespacedNameFrom(tc.in)
			assert.Equal(tc.out, actual)
			assert.Equal(tc.expectErr, err != nil)
		})
	}
}

func TestGetOSMControllerPods(t *testing.T) {
	testNamespace := "osm-namespace"

	tests := []struct {
		testName         string
		pods             []*corev1.Pod
		expectedPodNames []string
	}{
		{
			testName: "multiple pods (osm-controller pods and other pods) in multiple namespaces",
			pods: []*corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osm-controller-pod-1",
						Namespace: testNamespace,
						Labels: map[string]string{
							constants.AppLabel: constants.OSMControllerName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osm-controller-pod-2",
						Namespace: testNamespace,
						Labels: map[string]string{
							constants.AppLabel: constants.OSMControllerName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "some-other-osm-controller-pod",
						Namespace: "some-other-namespace",
						Labels: map[string]string{
							constants.AppLabel: constants.OSMControllerName,
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "application-pod-1",
						Namespace: testNamespace,
						Labels: map[string]string{
							constants.AppLabel: "myapp",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "application-pod-2",
						Namespace: "some-other-namespace",
						Labels: map[string]string{
							constants.AppLabel: "myapp",
						},
					},
				},
			},
			expectedPodNames: []string{
				"osm-controller-pod-1",
				"osm-controller-pod-2",
			},
		},
		{
			testName:         "no pods",
			pods:             []*corev1.Pod{},
			expectedPodNames: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			assert := tassert.New(t)

			objs := make([]runtime.Object, len(test.pods))
			for i := range test.pods {
				objs[i] = test.pods[i]
			}

			fakeClientSet := fake.NewSimpleClientset(objs...)
			podList := getOSMControllerPods(fakeClientSet, testNamespace)
			actualPodNames := make([]string, len(podList.Items))
			for i, pod := range podList.Items {
				actualPodNames[i] = pod.Name
			}
			assert.ElementsMatch(test.expectedPodNames, actualPodNames)
		})
	}
}
