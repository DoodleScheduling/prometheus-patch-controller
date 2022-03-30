<p>Packages:</p>
<ul>
<li>
<a href="#metrics.infra.doodle.com%2fv1beta1">metrics.infra.doodle.com/v1beta1</a>
</li>
</ul>
<h2 id="metrics.infra.doodle.com/v1beta1">metrics.infra.doodle.com/v1beta1</h2>
<div>
<p>Package v1beta1 contains API Schema definitions for the metrics.infra.doodle.com v1beta1 API group.</p>
</div>
Resource Types:
<ul></ul>
<h3 id="metrics.infra.doodle.com/v1beta1.JSON6902Patch">JSON6902Patch
</h3>
<p>
(<em>Appears on:</em><a href="#metrics.infra.doodle.com/v1beta1.PrometheusPatchRuleSpec">PrometheusPatchRuleSpec</a>)
</p>
<div>
<p>JSON6902Patch is a target selector and a list of JSON6902 patches</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>patch</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.JSONPatch">
[]JSONPatch
</a>
</em>
</td>
<td>
<p>Patch contains JSON6902 patches with
an array of operation objects.</p>
</td>
</tr>
<tr>
<td>
<code>target</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.Selector">
Selector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Target points to the resources that the patch document should be applied to.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metrics.infra.doodle.com/v1beta1.JSONPatch">JSONPatch
</h3>
<p>
(<em>Appears on:</em><a href="#metrics.infra.doodle.com/v1beta1.JSON6902Patch">JSON6902Patch</a>)
</p>
<div>
<p>JSONPatch is a JSON 6902 conform patch</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>op</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>path</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>value</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1?tab=doc#JSON">
Kubernetes pkg/apis/apiextensions/v1.JSON
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metrics.infra.doodle.com/v1beta1.PrometheusPatchRule">PrometheusPatchRule
</h3>
<div>
<p>PrometheusPatchRule is the Schema for the patchrules API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.PrometheusPatchRuleSpec">
PrometheusPatchRuleSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>prometheus</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.PrometheusSpec">
PrometheusSpec
</a>
</em>
</td>
<td>
<p>Prometheus holds information about where to find prometheus</p>
</td>
</tr>
<tr>
<td>
<code>interval</code><br/>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>Interval is the duration in which the expression gets evaluated</p>
</td>
</tr>
<tr>
<td>
<code>expr</code><br/>
<em>
string
</em>
</td>
<td>
<p>Expression is the prometheus .query</p>
</td>
</tr>
<tr>
<td>
<code>for</code><br/>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>For is a durstion for how long the rule should be in pending before apply patches.</p>
</td>
</tr>
<tr>
<td>
<code>json6902Patches</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.JSON6902Patch">
[]JSON6902Patch
</a>
</em>
</td>
<td>
<p>.JSON6902Patches define to what target are applied what patches</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend may suspend reconciliation of the resource.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.PrometheusPatchRuleStatus">
PrometheusPatchRuleStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metrics.infra.doodle.com/v1beta1.PrometheusPatchRuleSpec">PrometheusPatchRuleSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metrics.infra.doodle.com/v1beta1.PrometheusPatchRule">PrometheusPatchRule</a>)
</p>
<div>
<p>PrometheusPatchRuleSpec defines the desired state of PrometheusPatchRule</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>prometheus</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.PrometheusSpec">
PrometheusSpec
</a>
</em>
</td>
<td>
<p>Prometheus holds information about where to find prometheus</p>
</td>
</tr>
<tr>
<td>
<code>interval</code><br/>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>Interval is the duration in which the expression gets evaluated</p>
</td>
</tr>
<tr>
<td>
<code>expr</code><br/>
<em>
string
</em>
</td>
<td>
<p>Expression is the prometheus .query</p>
</td>
</tr>
<tr>
<td>
<code>for</code><br/>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>For is a durstion for how long the rule should be in pending before apply patches.</p>
</td>
</tr>
<tr>
<td>
<code>json6902Patches</code><br/>
<em>
<a href="#metrics.infra.doodle.com/v1beta1.JSON6902Patch">
[]JSON6902Patch
</a>
</em>
</td>
<td>
<p>.JSON6902Patches define to what target are applied what patches</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend may suspend reconciliation of the resource.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metrics.infra.doodle.com/v1beta1.PrometheusPatchRuleStatus">PrometheusPatchRuleStatus
</h3>
<p>
(<em>Appears on:</em><a href="#metrics.infra.doodle.com/v1beta1.PrometheusPatchRule">PrometheusPatchRule</a>)
</p>
<div>
<p>PrometheusPatchRuleStatus defines the observed state of PrometheusPatchRule</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions holds the conditions for the PrometheusPatchRule.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="metrics.infra.doodle.com/v1beta1.PrometheusSpec">PrometheusSpec
</h3>
<p>
(<em>Appears on:</em><a href="#metrics.infra.doodle.com/v1beta1.PrometheusPatchRuleSpec">PrometheusPatchRuleSpec</a>)
</p>
<div>
<p>PrometheusSpec contains specs for accessing prometheus</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>address</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="metrics.infra.doodle.com/v1beta1.Selector">Selector
</h3>
<p>
(<em>Appears on:</em><a href="#metrics.infra.doodle.com/v1beta1.JSON6902Patch">JSON6902Patch</a>)
</p>
<div>
<p>Selector specifies a set of resources. Any resource that matches intersection of all conditions is included in this
set.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>group</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Group is the API group to select resources from.
Together with Version and Kind it is capable of unambiguously identifying and/or selecting resources.
<a href="https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md">https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md</a></p>
</td>
</tr>
<tr>
<td>
<code>version</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Version of the API Group to select resources from.
Together with Group and Kind it is capable of unambiguously identifying and/or selecting resources.
<a href="https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md">https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md</a></p>
</td>
</tr>
<tr>
<td>
<code>resource</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Kind of the API Group to select resources from.
Together with Group and Version it is capable of unambiguously
identifying and/or selecting resources.
<a href="https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md">https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md</a></p>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Namespace to select resources from.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name to match resources with.</p>
</td>
</tr>
<tr>
<td>
<code>labelSelector</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>LabelSelector is a string that follows the label selection expression
<a href="https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#api">https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#api</a>
It matches with the resource labels.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
.
</em></p>
