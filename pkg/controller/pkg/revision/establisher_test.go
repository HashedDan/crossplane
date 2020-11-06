/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package revision

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/crossplane/crossplane/apis/pkg/v1alpha1"
)

var _ Establisher = &APIEstablisher{}

func TestAPIEstablisherEstablish(t *testing.T) {
	errBoom := errors.New("boom")
	trueVal := true

	type args struct {
		est     *APIEstablisher
		objs    []runtime.Object
		parent  resource.Object
		control bool
	}

	type want struct {
		err  error
		refs []runtimev1alpha1.TypedReference
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"SuccessfulExistsEstablishControl": {
			reason: "Establishment should be successful if we can establish control for a parent of existing objects.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet:    test.NewMockGetFn(nil),
						MockUpdate: test.NewMockUpdateFn(nil),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent:  &v1alpha1.ProviderRevision{},
				control: true,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulExistsAlreadyControlled": {
			reason: "Establishment should be successful if objects exist and are already controlled.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet: test.NewMockGetFn(nil, func(o runtime.Object) error {
							crd := &apiextensions.CustomResourceDefinition{}
							crd.SetName("ref-me")
							crd.SetOwnerReferences([]metav1.OwnerReference{
								{
									APIVersion: v1alpha1.SchemeGroupVersion.String(),
									Kind:       v1alpha1.ProviderRevisionKind,
									Name:       "test-pr",
									Controller: &trueVal,
								},
							})
							*o.(*apiextensions.CustomResourceDefinition) = *crd
							return nil
						}),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent: &v1alpha1.ProviderRevision{
					TypeMeta: metav1.TypeMeta{
						APIVersion: v1alpha1.SchemeGroupVersion.String(),
						Kind:       v1alpha1.ProviderRevisionKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-pr",
					},
				},
				control: true,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulExistsOwnToControl": {
			reason: "Establishment should be successful if an owned object is successfully transitioned to controlled.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet: test.NewMockGetFn(nil, func(o runtime.Object) error {
							crd := &apiextensions.CustomResourceDefinition{}
							crd.SetName("ref-me")
							crd.SetOwnerReferences([]metav1.OwnerReference{
								{
									APIVersion: v1alpha1.SchemeGroupVersion.String(),
									Kind:       v1alpha1.ProviderRevisionKind,
									Name:       "test-pr",
								},
							})
							*o.(*apiextensions.CustomResourceDefinition) = *crd
							return nil
						}),
						MockUpdate: test.NewMockUpdateFn(nil),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent: &v1alpha1.ProviderRevision{
					TypeMeta: metav1.TypeMeta{
						APIVersion: v1alpha1.SchemeGroupVersion.String(),
						Kind:       v1alpha1.ProviderRevisionKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-pr",
					},
				},
				control: true,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulNotExistsEstablishControl": {
			reason: "Establishment should be successful if we can establish control for a parent of new objects.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet:    test.NewMockGetFn(kerrors.NewNotFound(schema.GroupResource{}, "")),
						MockCreate: test.NewMockCreateFn(nil),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent:  &v1alpha1.ProviderRevision{},
				control: true,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulExistsEstablishOwnership": {
			reason: "Establishment should be successful if we can establish ownership for a parent of existing objects.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet:    test.NewMockGetFn(nil),
						MockUpdate: test.NewMockUpdateFn(nil),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent:  &v1alpha1.ProviderRevision{},
				control: false,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulExistsAlreadyOwned": {
			reason: "Establishment should be successful if objects exist and are already owned.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet: test.NewMockGetFn(nil, func(o runtime.Object) error {
							crd := &apiextensions.CustomResourceDefinition{}
							crd.SetName("ref-me")
							crd.SetOwnerReferences([]metav1.OwnerReference{
								{
									APIVersion: v1alpha1.SchemeGroupVersion.String(),
									Kind:       v1alpha1.ProviderRevisionKind,
									Name:       "test-pr",
								},
							})
							*o.(*apiextensions.CustomResourceDefinition) = *crd
							return nil
						}),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent: &v1alpha1.ProviderRevision{
					TypeMeta: metav1.TypeMeta{
						APIVersion: v1alpha1.SchemeGroupVersion.String(),
						Kind:       v1alpha1.ProviderRevisionKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-pr",
					},
				},
				control: false,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulExistsControlToOwn": {
			reason: "Establishment should be successful if a controlled object is successfully transitioned to owned.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet: test.NewMockGetFn(nil, func(o runtime.Object) error {
							crd := &apiextensions.CustomResourceDefinition{}
							crd.SetName("ref-me")
							crd.SetOwnerReferences([]metav1.OwnerReference{
								{
									APIVersion: v1alpha1.SchemeGroupVersion.String(),
									Kind:       v1alpha1.ProviderRevisionKind,
									Name:       "test-pr",
									Controller: &trueVal,
								},
							})
							*o.(*apiextensions.CustomResourceDefinition) = *crd
							return nil
						}),
						MockUpdate: test.NewMockUpdateFn(nil),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent: &v1alpha1.ProviderRevision{
					TypeMeta: metav1.TypeMeta{
						APIVersion: v1alpha1.SchemeGroupVersion.String(),
						Kind:       v1alpha1.ProviderRevisionKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-pr",
					},
				},
				control: false,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"SuccessfulNotExistsEstablishOwnership": {
			reason: "Establishment should be successful if we can establish ownership for a parent of new objects.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet:    test.NewMockGetFn(kerrors.NewNotFound(schema.GroupResource{}, "")),
						MockCreate: test.NewMockCreateFn(nil),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent:  &v1alpha1.ProviderRevision{},
				control: false,
			},
			want: want{
				refs: []runtimev1alpha1.TypedReference{{Name: "ref-me"}},
			},
		},
		"FailedCreate": {
			reason: "Cannot establish control of object if we cannot create it.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet:    test.NewMockGetFn(kerrors.NewNotFound(schema.GroupResource{}, "")),
						MockCreate: test.NewMockCreateFn(errBoom),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent: &v1alpha1.ProviderRevision{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
				},
				control: true,
			},
			want: want{
				err: errBoom,
			},
		},
		"FailedUpdate": {
			reason: "Cannot establish control of object if we cannot update it.",
			args: args{
				est: &APIEstablisher{
					client: &test.MockClient{
						MockGet:    test.NewMockGetFn(nil),
						MockUpdate: test.NewMockUpdateFn(errBoom),
					},
				},
				objs: []runtime.Object{
					&apiextensions.CustomResourceDefinition{
						ObjectMeta: metav1.ObjectMeta{
							Name: "ref-me",
						},
					},
				},
				parent: &v1alpha1.ProviderRevision{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
					},
				},
				control: true,
			},
			want: want{
				err: errBoom,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			refs, err := tc.args.est.Establish(context.TODO(), tc.args.objs, tc.args.parent, tc.args.control)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Check(...): -want error, +got error:\n%s", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.refs, refs, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Check(...): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}
