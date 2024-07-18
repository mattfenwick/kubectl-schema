package swagger

import (
	"fmt"
	"github.com/mattfenwick/collections/pkg/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func RunShowResourcesTests() {
	versions := []string{"1.18.20", "1.20.15", "1.22.12", "1.24.0", "1.25.0-alpha.3"}
	resources := set.FromSlice([]string{"Ingress", "CronJob", "CustomResourceDefinition"})
	include := func(apiVersion string, resource string) bool {
		return resources.Contains(resource)
	}

	Describe("Show resource", func() {
		It("By resource -- no diff", func() {
			actual := ShowResources(ShowResourcesGroupByResource, versions, include, false, ShowResourcesFormatTable)
			Expect(actual).To(Equal(byResourceNoDiff[1:]))
		})
		It("By apiversion -- no diff", func() {
			actual := ShowResources(ShowResourcesGroupByApiVersion, versions, include, false, ShowResourcesFormatTable)
			Expect(actual).To(Equal(byApiVersionNoDiff[1:]))
		})
		It("By resource -- diff", func() {
			actual := ShowResources(ShowResourcesGroupByResource, versions, include, true, ShowResourcesFormatTable)
			fmt.Printf("expect:\n%s\n", byResourceWithDiff[1:])
			fmt.Printf("actual:\n%s\n", actual)
			Expect(actual).To(Equal(byResourceWithDiff[1:]))
		})
		It("By apiversion -- diff", func() {
			actual := ShowResources(ShowResourcesGroupByApiVersion, versions, include, true, ShowResourcesFormatTable)
			fmt.Printf("actual vs. expected:\n%s\n\n%s\n\n", actual, byApiVersionWithDiff)
			Expect(actual).To(Equal(byApiVersionWithDiff[1:]))
		})
	})
}

var (
	byResourceNoDiff = `
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
|         RESOURCE         |           1.18.20            |           1.20.15            |         1.22.12         |         1.24.0          |     1.25.0-ALPHA 3      |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
| CronJob                  | batch.v1beta1                | batch.v1beta1                | batch.v1                | batch.v1                | batch.v1                |
|                          | batch.v2alpha1               | batch.v2alpha1               | batch.v1beta1           | batch.v1beta1           |                         |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
| CustomResourceDefinition | apiextensions.k8s.io.v1      | apiextensions.k8s.io.v1      | apiextensions.k8s.io.v1 | apiextensions.k8s.io.v1 | apiextensions.k8s.io.v1 |
|                          | apiextensions.k8s.io.v1beta1 | apiextensions.k8s.io.v1beta1 |                         |                         |                         |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
| Ingress                  | extensions.v1beta1           | extensions.v1beta1           | networking.k8s.io.v1    | networking.k8s.io.v1    | networking.k8s.io.v1    |
|                          | networking.k8s.io.v1beta1    | networking.k8s.io.v1         |                         |                         |                         |
|                          |                              | networking.k8s.io.v1beta1    |                         |                         |                         |
+--------------------------+------------------------------+------------------------------+-------------------------+-------------------------+-------------------------+
`

	byApiVersionNoDiff = `
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
|         API VERSION          |         1.18.20          |         1.20.15          |         1.22.12          |          1.24.0          |      1.25.0-ALPHA 3      |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
| apiextensions.k8s.io.v1      | CustomResourceDefinition | CustomResourceDefinition | CustomResourceDefinition | CustomResourceDefinition | CustomResourceDefinition |
+------------------------------+                          +                          +--------------------------+--------------------------+--------------------------+
| apiextensions.k8s.io.v1beta1 |                          |                          |                          |                          |                          |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
| batch.v1                     |                          |                          | CronJob                  | CronJob                  | CronJob                  |
+------------------------------+--------------------------+--------------------------+                          +                          +--------------------------+
| batch.v1beta1                | CronJob                  | CronJob                  |                          |                          |                          |
+------------------------------+                          +                          +--------------------------+--------------------------+--------------------------+
| batch.v2alpha1               |                          |                          |                          |                          |                          |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
| extensions.v1beta1           | Ingress                  | Ingress                  |                          |                          |                          |
+------------------------------+--------------------------+                          +--------------------------+--------------------------+--------------------------+
| networking.k8s.io.v1         |                          |                          | Ingress                  | Ingress                  | Ingress                  |
+------------------------------+--------------------------+                          +--------------------------+--------------------------+--------------------------+
| networking.k8s.io.v1beta1    | Ingress                  |                          |                          |                          |                          |
+------------------------------+--------------------------+--------------------------+--------------------------+--------------------------+--------------------------+
`

	byResourceWithDiff = `
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
|         RESOURCE         |           1.18.20            |        1.20.15         |            1.22.12             | 1.24.0 | 1.25.0-ALPHA 3  |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
| CronJob                  | batch.v1beta1                |                        | add:                           |        | remove:         |
|                          | batch.v2alpha1               |                        |   batch.v1                     |        |   batch.v1beta1 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        | remove:                        |        |                 |
|                          |                              |                        |   batch.v2alpha1               |        |                 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        |                                |        |                 |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
| CustomResourceDefinition | apiextensions.k8s.io.v1      |                        | remove:                        |        |                 |
|                          | apiextensions.k8s.io.v1beta1 |                        |   apiextensions.k8s.io.v1beta1 |        |                 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        |                                |        |                 |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
| Ingress                  | extensions.v1beta1           | add:                   | remove:                        |        |                 |
|                          | networking.k8s.io.v1beta1    |   networking.k8s.io.v1 |   extensions.v1beta1           |        |                 |
|                          |                              |                        |   networking.k8s.io.v1beta1    |        |                 |
|                          |                              |                        |                                |        |                 |
|                          |                              |                        |                                |        |                 |
+--------------------------+------------------------------+------------------------+--------------------------------+--------+-----------------+
`

	byApiVersionWithDiff = `
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
|         API VERSION          |         1.18.20          |  1.20.15  |          1.22.12           | 1.24.0 | 1.25.0-ALPHA 3 |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| apiextensions.k8s.io.v1      | CustomResourceDefinition |           |                            |        |                |
+------------------------------+                          +-----------+----------------------------+--------+----------------+
| apiextensions.k8s.io.v1beta1 |                          |           | remove:                    |        |                |
|                              |                          |           |   CustomResourceDefinition |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| batch.v1                     |                          |           | add:                       |        |                |
|                              |                          |           |   CronJob                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| batch.v1beta1                | CronJob                  |           |                            |        | remove:        |
|                              |                          |           |                            |        |   CronJob      |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+                          +-----------+----------------------------+--------+----------------+
| batch.v2alpha1               |                          |           | remove:                    |        |                |
|                              |                          |           |   CronJob                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| extensions.v1beta1           | Ingress                  |           | remove:                    |        |                |
|                              |                          |           |   Ingress                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| networking.k8s.io.v1         |                          | add:      |                            |        |                |
|                              |                          |   Ingress |                            |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
| networking.k8s.io.v1beta1    | Ingress                  |           | remove:                    |        |                |
|                              |                          |           |   Ingress                  |        |                |
|                              |                          |           |                            |        |                |
|                              |                          |           |                            |        |                |
+------------------------------+--------------------------+-----------+----------------------------+--------+----------------+
`
)
