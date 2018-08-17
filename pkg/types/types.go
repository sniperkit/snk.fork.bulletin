/*
Sniperkit-Bot
- Status: analyzed
*/

package types

import (
	yaml "gopkg.in/yaml.v2"
)

type Comparator struct {
	s1  []string
	s2  []string
	s1a [][]string
	s2a [][]string
	b1  []bool
	b2  []bool
	sm  []map[string]string
}

func (c *Comparator) Bool(a, b bool) *Comparator {
	c.b1 = append(c.b1, a)
	c.b2 = append(c.b2, b)
	return c
}

func (c *Comparator) Strings(a, b string) *Comparator {
	c.s1 = append(c.s1, a)
	c.s2 = append(c.s2, b)
	return c
}

func (c *Comparator) StringSlice(a, b []string) *Comparator {
	c.s1a = append(c.s1a, a)
	c.s2a = append(c.s2a, b)
	return c
}

func (c *Comparator) Equal() bool {
	if !StringSliceEqual(c.s1, c.s2) {
		return false
	}
	for i := 0; i < len(c.s1a); i++ {
		if !StringSliceEqual(c.s1a[i], c.s2a[i]) {
			return false
		}
	}
	for i := 0; i < len(c.b1); i++ {
		if c.b1[i] != c.b2[i] {
			return false
		}
	}
	return true
}

func StringMapEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		if v != bv {
			return false
		}
	}
	return true
}

func GetStringMap(i interface{}) (map[string]string, error) {
	res := make(map[string]string)
	d, err := yaml.Marshal(i)
	if err != nil {
		return res, err
	}
	err = yaml.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		bv := b[i]
		if v != bv {
			return false
		}
	}
	return true
}

type Comparable interface {
	Equal(interface{}) bool
}

type Set interface {
	Get() []Comparable
	Add(Comparable)
}

// don't complain, a very inefficient implementation
type SimpleSet struct {
	ca []Comparable
}

func NewSet(eles ...Comparable) *SimpleSet {
	res := &SimpleSet{}
	for _, ele := range eles {
		res.Add(ele)
	}
	return res
}

func (o *SimpleSet) Get() []Comparable {
	return o.ca
}

func (o *SimpleSet) Add(t Comparable) {
	found := false
	for _, r := range o.ca {
		if r.Equal(t) {
			found = true
			break
		}
	}
	if !found {
		o.ca = append(o.ca, t)
	}
}

type InternalError string

func (i InternalError) Error() string {
	return string(i)
}
